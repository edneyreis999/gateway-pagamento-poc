import { NotFoundError } from '../../../../../shared/domain/errors/not-found.error';
import { Uuid } from '../../../../../shared/domain/value-objects/uuid.vo';
import { setupSequelize } from '../../../../../shared/infra/testing/helpers';
import { User } from '../../../../domain/user.entity';
import { UserSearchParams, UserSearchResult } from '../../../../domain/user.repository';
import { UserModelMapper } from './user-model-mapper';
import { UserSequelizeRepository } from './user-sequelize.repository';
import { UserModel } from './user.model';

describe('UserSequelizeRepository Integration Test', () => {
  let repository: UserSequelizeRepository;
  setupSequelize({ models: [UserModel] });

  beforeEach(async () => {
    repository = new UserSequelizeRepository(UserModel);
  });

  it('should inserts a new entity', async () => {
    const user = User.fake().aUser().build();
    await repository.insert(user);
    const entity = await repository.findById(user.userId);
    expect(entity!.toJSON()).toStrictEqual(user.toJSON());
  });

  it('should finds a entity by id', async () => {
    let entityFound = await repository.findById(new Uuid());
    expect(entityFound).toBeNull();

    const entity = User.fake().aUser().build();
    await repository.insert(entity);
    entityFound = await repository.findById(entity.userId);
    expect(entity.toJSON()).toStrictEqual(entityFound!.toJSON());
  });

  it('should return all users', async () => {
    const entity = User.fake().aUser().build();
    await repository.insert(entity);
    const entities = await repository.findAll();
    expect(entities).toHaveLength(1);
    expect(JSON.stringify(entities)).toBe(JSON.stringify([entity]));
  });

  it('should throw error on update when a entity not found', async () => {
    const entity = User.fake().aUser().build();
    await expect(repository.update(entity)).rejects.toThrow(
      new NotFoundError(entity.userId.id, User),
    );
  });

  it('should update a entity', async () => {
    const entity = User.fake().aUser().build();
    await repository.insert(entity);

    entity.changeDisplayName('Movie updated');
    await repository.update(entity);

    const entityFound = await repository.findById(entity.userId);
    expect(entity.toJSON()).toStrictEqual(entityFound!.toJSON());
  });

  it('should throw error on delete when a entity not found', async () => {
    const userId = new Uuid();
    await expect(repository.delete(userId)).rejects.toThrow(new NotFoundError(userId.id, User));
  });

  it('should delete a entity', async () => {
    const entity = new User({ displayName: 'Movie' });
    await repository.insert(entity);

    await repository.delete(entity.userId);
    await expect(repository.findById(entity.userId)).resolves.toBeNull();
  });

  describe('search method tests', () => {
    it('should only apply paginate when other params are null', async () => {
      const createdAt = new Date();
      const users = User.fake()
        .theUsers(16)
        .withDisplayName('Movie')
        .withDustBalance(null)
        .withCreatedAt(createdAt)
        .build();
      await repository.bulkInsert(users);
      const spyToEntity = jest.spyOn(UserModelMapper, 'toEntity');

      const searchOutput = await repository.search(new UserSearchParams());
      expect(searchOutput).toBeInstanceOf(UserSearchResult);
      expect(spyToEntity).toHaveBeenCalledTimes(15);
      expect(searchOutput.toJSON()).toMatchObject({
        total: 16,
        current_page: 1,
        last_page: 2,
        per_page: 15,
      });
      searchOutput.items.forEach(item => {
        expect(item).toBeInstanceOf(User);
        expect(item.userId).toBeDefined();
      });
      const items = searchOutput.items.map(item => item.toJSON());
      expect(items).toMatchObject(
        new Array(15).fill({
          displayName: 'Movie',
          dustBalance: 0,
          isActive: true,
          createdAt: createdAt,
        }),
      );
    });

    it('should order by createdAt DESC when search params are null', async () => {
      const createdAt = new Date();
      const users = User.fake()
        .theUsers(16)
        .withDisplayName(index => `Movie ${index}`)
        .withDustBalance(null)
        .withCreatedAt(index => new Date(createdAt.getTime() + index))
        .build();
      const searchOutput = await repository.search(new UserSearchParams());
      const items = searchOutput.items;
      [...items].reverse().forEach((item, index) => {
        expect(`Movie ${index}`).toBe(`${users[index + 1].displayName}`);
      });
    });

    it('should apply paginate and filter', async () => {
      const users = [
        User.fake()
          .aUser()
          .withDisplayName('test')
          .withCreatedAt(new Date(new Date().getTime() + 5000))
          .build(),
        User.fake()
          .aUser()
          .withDisplayName('a')
          .withCreatedAt(new Date(new Date().getTime() + 4000))
          .build(),
        User.fake()
          .aUser()
          .withDisplayName('TEST')
          .withCreatedAt(new Date(new Date().getTime() + 3000))
          .build(),
        User.fake()
          .aUser()
          .withDisplayName('TeSt')
          .withCreatedAt(new Date(new Date().getTime() + 1000))
          .build(),
      ];

      await repository.bulkInsert(users);

      let searchOutput = await repository.search(
        new UserSearchParams({
          page: 1,
          per_page: 2,
          filter: 'TEST',
        }),
      );
      expect(searchOutput.toJSON(true)).toMatchObject(
        new UserSearchResult({
          items: [users[0], users[2]],
          total: 3,
          current_page: 1,
          per_page: 2,
        }).toJSON(true),
      );

      searchOutput = await repository.search(
        new UserSearchParams({
          page: 2,
          per_page: 2,
          filter: 'TEST',
        }),
      );
      expect(searchOutput.toJSON(true)).toMatchObject(
        new UserSearchResult({
          items: [users[3]],
          total: 3,
          current_page: 2,
          per_page: 2,
        }).toJSON(true),
      );
    });

    it('should apply paginate and sort', async () => {
      expect(repository.sortableFields).toStrictEqual(['displayName', 'createdAt']);

      const users = [
        User.fake().aUser().withDisplayName('b').build(),
        User.fake().aUser().withDisplayName('a').build(),
        User.fake().aUser().withDisplayName('d').build(),
        User.fake().aUser().withDisplayName('e').build(),
        User.fake().aUser().withDisplayName('c').build(),
      ];
      await repository.bulkInsert(users);

      const arrange = [
        {
          params: new UserSearchParams({
            page: 1,
            per_page: 2,
            sort: 'displayName',
          }),
          result: new UserSearchResult({
            items: [users[1], users[0]],
            total: 5,
            current_page: 1,
            per_page: 2,
          }),
        },
        {
          params: new UserSearchParams({
            page: 2,
            per_page: 2,
            sort: 'displayName',
          }),
          result: new UserSearchResult({
            items: [users[4], users[2]],
            total: 5,
            current_page: 2,
            per_page: 2,
          }),
        },
        {
          params: new UserSearchParams({
            page: 1,
            per_page: 2,
            sort: 'displayName',
            sort_dir: 'desc',
          }),
          result: new UserSearchResult({
            items: [users[3], users[2]],
            total: 5,
            current_page: 1,
            per_page: 2,
          }),
        },
        {
          params: new UserSearchParams({
            page: 2,
            per_page: 2,
            sort: 'displayName',
            sort_dir: 'desc',
          }),
          result: new UserSearchResult({
            items: [users[4], users[0]],
            total: 5,
            current_page: 2,
            per_page: 2,
          }),
        },
      ];

      for (const i of arrange) {
        const result = await repository.search(i.params);
        expect(result.toJSON(true)).toMatchObject(i.result.toJSON(true));
      }
    });

    describe('should search using filter, sort and paginate', () => {
      const users = [
        User.fake().aUser().withDisplayName('test').build(),
        User.fake().aUser().withDisplayName('a').build(),
        User.fake().aUser().withDisplayName('TEST').build(),
        User.fake().aUser().withDisplayName('e').build(),
        User.fake().aUser().withDisplayName('TeSt').build(),
      ];

      const arrange = [
        {
          search_params: new UserSearchParams({
            page: 1,
            per_page: 2,
            sort: 'displayName',
            filter: 'TEST',
          }),
          search_result: new UserSearchResult({
            items: [users[2], users[4]],
            total: 3,
            current_page: 1,
            per_page: 2,
          }),
        },
        {
          search_params: new UserSearchParams({
            page: 2,
            per_page: 2,
            sort: 'displayName',
            filter: 'TEST',
          }),
          search_result: new UserSearchResult({
            items: [users[0]],
            total: 3,
            current_page: 2,
            per_page: 2,
          }),
        },
      ];

      beforeEach(async () => {
        await repository.bulkInsert(users);
      });

      test.each(arrange)(
        'when value is $search_params',
        async ({ search_params, search_result }) => {
          const result = await repository.search(search_params);
          expect(result.toJSON(true)).toMatchObject(search_result.toJSON(true));
        },
      );
    });
  });
});
