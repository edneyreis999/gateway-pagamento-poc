import { User } from '../../domain/user.entity';
import { UserInMemoryRepository } from './user-in-memory.repository';

describe('UserInMemoryRepository', () => {
  let repository: UserInMemoryRepository;

  beforeEach(() => (repository = new UserInMemoryRepository()));
  it('should no filter items when filter object is null', async () => {
    const items = [User.fake().aUser().build()];
    const filterSpy = jest.spyOn(items, 'filter' as any);

    const itemsFiltered = await repository['applyFilter'](items, null);
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

  it('should sort by created_at when sort param is null', async () => {
    const created_at = new Date();

    const items = [
      User.fake().aUser().withDisplayName('test').withCreatedAt(created_at).build(),
      User.fake()
        .aUser()
        .withDisplayName('TEST')
        .withCreatedAt(new Date(created_at.getTime() + 100))
        .build(),
      User.fake()
        .aUser()
        .withDisplayName('fake')
        .withCreatedAt(new Date(created_at.getTime() + 200))
        .build(),
    ];

    const itemsSorted = await repository['applySort'](items, null, null);
    expect(itemsSorted).toStrictEqual([items[2], items[1], items[0]]);
  });

  it('should sort by name', async () => {
    const items = [
      User.fake().aUser().withDisplayName('c').build(),
      User.fake().aUser().withDisplayName('b').build(),
      User.fake().aUser().withDisplayName('a').build(),
    ];

    let itemsSorted = await repository['applySort'](items, 'displayName', 'asc');
    expect(itemsSorted).toStrictEqual([items[2], items[1], items[0]]);

    itemsSorted = await repository['applySort'](items, 'displayName', 'desc');
    expect(itemsSorted).toStrictEqual([items[0], items[1], items[2]]);
  });
});
