import { setupSequelize } from '../../../../shared/infra/testing/helpers';
import { User } from '../../../domain/user.entity';
import { UserSequelizeRepository } from '../../../infra/db/sequelize/sequelize/user-sequelize.repository';
import { UserModel } from '../../../infra/db/sequelize/sequelize/user.model';
import { UserOutputMapper } from '../_user-shared/user-output';
import { ListUsersUseCase } from './list-users.use-case';

describe('ListUsersUseCase Integration Tests', () => {
  let useCase: ListUsersUseCase;
  let repository: UserSequelizeRepository;

  setupSequelize({ models: [UserModel] });

  beforeEach(() => {
    repository = new UserSequelizeRepository(UserModel);
    useCase = new ListUsersUseCase(repository);
  });

  it('should return output sorted by created_at when input param is empty', async () => {
    const users = User.fake()
      .theUsers(2)
      .withCreatedAt(i => new Date(new Date().getTime() + 1000 + i))
      .build();

    await repository.bulkInsert(users);
    const output = await useCase.execute({});
    console.log(output);
    expect(output).toEqual({
      items: [...users].reverse().map(UserOutputMapper.toOutput),
      total: 2,
      current_page: 1,
      per_page: 15,
      last_page: 1,
    });
  });

  it('should returns output using pagination, sort and filter', async () => {
    const users = [
      new User({ displayName: 'a' }),
      new User({
        displayName: 'AAA',
      }),
      new User({
        displayName: 'AaA',
      }),
      new User({
        displayName: 'b',
      }),
      new User({
        displayName: 'c',
      }),
    ];
    await repository.bulkInsert(users);
    const spyUpdate = jest.spyOn(repository, 'search');

    let output = await useCase.execute({
      page: 1,
      per_page: 2,
      sort: 'displayName',
      filter: 'a',
    });
    expect(spyUpdate).toHaveBeenCalledTimes(1);
    expect(output).toEqual({
      items: [users[1], users[2]].map(UserOutputMapper.toOutput),
      total: 3,
      current_page: 1,
      per_page: 2,
      last_page: 2,
    });

    output = await useCase.execute({
      page: 2,
      per_page: 2,
      sort: 'displayName',
      filter: 'a',
    });
    expect(output).toEqual({
      items: [users[0]].map(UserOutputMapper.toOutput),
      total: 3,
      current_page: 2,
      per_page: 2,
      last_page: 2,
    });

    output = await useCase.execute({
      page: 1,
      per_page: 2,
      sort: 'displayName',
      sort_dir: 'desc',
      filter: 'a',
    });
    expect(output).toEqual({
      items: [users[0], users[2]].map(UserOutputMapper.toOutput),
      total: 3,
      current_page: 1,
      per_page: 2,
      last_page: 2,
    });
  });
});
