import { User } from '../../../domain/user.entity';
import { UserInMemoryRepository } from './user-in-memory.repository';

describe('UserInMemoryRepository', () => {
  let repository: UserInMemoryRepository;

  beforeEach(() => (repository = new UserInMemoryRepository()));
  it('should no filter items when filter object is null', async () => {
    const items = [User.fake().aUser().build()];
    const filterSpy = jest.spyOn(items, 'filter' as any);

    const itemsFiltered = await repository['applyFilter'](items, null as any);
    expect(filterSpy).not.toHaveBeenCalled();
    expect(itemsFiltered).toStrictEqual(items);
  });

  it('should filter items using filter parameter', async () => {
    const items = [
      User.fake().aUser().withDisplayName('test').build(),
      User.fake().aUser().withDisplayName('TEST').build(),
      User.fake().aUser().withDisplayName('fake').build(),
    ];
    const filterSpy = jest.spyOn(items, 'filter' as any);

    const itemsFiltered = await repository['applyFilter'](items, 'TEST');
    expect(filterSpy).toHaveBeenCalledTimes(1);
    expect(itemsFiltered).toStrictEqual([items[0], items[1]]);
  });

  it('should sort by createdAt when sort param is null', async () => {
    const createdAt = new Date();

    const items = [
      User.fake().aUser().withDisplayName('test').withCreatedAt(createdAt).build(),
      User.fake()
        .aUser()
        .withDisplayName('TEST')
        .withCreatedAt(new Date(createdAt.getTime() + 100))
        .build(),
      User.fake()
        .aUser()
        .withDisplayName('fake')
        .withCreatedAt(new Date(createdAt.getTime() + 200))
        .build(),
    ];

    const itemsSorted = repository['applySort'](items, null, null);
    expect(itemsSorted).toStrictEqual([items[2], items[1], items[0]]);
  });

  it('should sort by displayName', async () => {
    const items = [
      User.fake().aUser().withDisplayName('c').build(),
      User.fake().aUser().withDisplayName('b').build(),
      User.fake().aUser().withDisplayName('a').build(),
    ];

    let itemsSorted = repository['applySort'](items, 'displayName', 'asc');
    expect(itemsSorted).toStrictEqual([items[2], items[1], items[0]]);

    itemsSorted = repository['applySort'](items, 'displayName', 'desc');
    expect(itemsSorted).toStrictEqual([items[0], items[1], items[2]]);
  });
});
