import {
  PaginationOutput,
  PaginationOutputMapper,
} from '../../../../shared/application/pagination-output';
import { IUseCase } from '../../../../shared/application/use-case.interface';
import {
  IUserRepository,
  UserSearchParams,
  UserSearchResult,
} from '../../../domain/user.repository';
import { UserOutputMapper, type UserOutput } from '../_user-shared/user-output';
import { ValidateListUsersInput, type ListUsersInput } from './list-users.use-case.input';

export class ListUsersUseCase implements IUseCase<ListUsersInput, ListUsersOutput> {
  constructor(private userRepo: IUserRepository) {}

  async execute(input: ListUsersInput): Promise<ListUsersOutput> {
    const errors = ValidateListUsersInput.validate(input);
    if (errors.length > 0) {
      throw errors;
    }

    const params = new UserSearchParams(input);
    const searchResult = await this.userRepo.search(params);
    return this.toOutput(searchResult);
  }

  private toOutput(searchResult: UserSearchResult): ListUsersOutput {
    const { items: _items } = searchResult;
    const items = _items.map(i => {
      return UserOutputMapper.toOutput(i);
    });
    return PaginationOutputMapper.toOutput(items, searchResult);
  }
}

export type ListUsersOutput = PaginationOutput<UserOutput>;
