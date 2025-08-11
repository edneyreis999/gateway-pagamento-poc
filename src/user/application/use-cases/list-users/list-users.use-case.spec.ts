import { User } from '../../../domain/user.entity';
import { UserSearchResult } from '../../../domain/user.repository';
import { UserInMemoryRepository } from '../../../infra/db/in-memory/user-in-memory.repository';
import { UserOutputMapper } from '../_user-shared/user-output';
import { ListUsersUseCase } from './list-users.use-case';

describe('ListUsersUseCase Unit Tests', () => {
  let listUsersUseCase: ListUsersUseCase;
  let repository: UserInMemoryRepository;

  beforeEach(() => {
    repository = new UserInMemoryRepository();
    listUsersUseCase = new ListUsersUseCase(repository);
  });

  test('toOutput method', () => {
    let result = new UserSearchResult({
      items: [],
      total: 1,
      current_page: 1,
      per_page: 2,
    });
    let output = listUsersUseCase['toOutput'](result);
    expect(output).toStrictEqual({
      items: [],
      total: 1,
      current_page: 1,
      per_page: 2,
      last_page: 1,
    });

    const entity = User.create({ displayName: 'Movie' });
    result = new UserSearchResult({
      items: [entity],
      total: 1,
      current_page: 1,
      per_page: 2,
    });

    output = listUsersUseCase['toOutput'](result);
    expect(output).toStrictEqual({
      items: [entity].map(UserOutputMapper.toOutput),
      total: 1,
      current_page: 1,
      per_page: 2,
      last_page: 1,
    });
  });

  it('should return output sorted by createdAt when input param is empty', async () => {
    const items = [
      new User({ displayName: 'test 1' }),
      new User({
        displayName: 'test 2',
        createdAt: new Date(new Date().getTime() + 100),
      }),
    ];
    repository.items = items;

    const output = await listUsersUseCase.execute({});
    expect(output).toStrictEqual({
      items: [...items].reverse().map(UserOutputMapper.toOutput),
      total: 2,
      current_page: 1,
      per_page: 15,
      last_page: 1,
    });
  });

  it('should return output using pagination, sort and filter', async () => {
    const items = [
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
    repository.items = items;

    let output = await listUsersUseCase.execute({
      page: 1,
      per_page: 2,
      sort: 'displayName',
      filter: 'a',
    });
    expect(output).toStrictEqual({
      items: [items[1], items[2]].map(UserOutputMapper.toOutput),
      total: 3,
      current_page: 1,
      per_page: 2,
      last_page: 2,
    });

    output = await listUsersUseCase.execute({
      page: 2,
      per_page: 2,
      sort: 'displayName',
      filter: 'a',
    });
    expect(output).toStrictEqual({
      items: [items[0]].map(UserOutputMapper.toOutput),
      total: 3,
      current_page: 2,
      per_page: 2,
      last_page: 2,
    });

    output = await listUsersUseCase.execute({
      page: 1,
      per_page: 2,
      sort: 'displayName',
      sort_dir: 'desc',
      filter: 'a',
    });
    expect(output).toStrictEqual({
      items: [items[0], items[2]].map(UserOutputMapper.toOutput),
      total: 3,
      current_page: 1,
      per_page: 2,
      last_page: 2,
    });
  });
});
