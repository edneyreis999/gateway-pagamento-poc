import { IsEnum, IsNumber, IsOptional, validateSync } from 'class-validator';
import type { SortDirection } from '../../../../shared/domain/repository/search-params';
import type { UserFilter } from '../../../domain/user.repository';

export type ListUsersInputConstructorProps = {
  page?: number;
  per_page?: number;
  sort?: 'displayName' | 'createdAt' | null;
  sort_dir?: SortDirection | null;
  filter?: UserFilter | null;
};

export class ListUsersInput {
  @IsOptional()
  @IsNumber()
  page?: number;

  @IsOptional()
  @IsNumber()
  per_page?: number;

  @IsOptional()
  @IsEnum(['displayName', 'createdAt'])
  sort?: 'displayName' | 'createdAt' | null;

  @IsOptional()
  @IsEnum(['asc', 'desc'])
  sort_dir?: SortDirection | null;

  @IsOptional()
  filter?: UserFilter | null;

  constructor(props?: ListUsersInputConstructorProps) {
    if (!props) return;

    props.page && (this.page = props.page);
    props.per_page && (this.per_page = props.per_page);
    props.sort && (this.sort = props.sort);
    props.sort_dir && (this.sort_dir = props.sort_dir);
    props.filter && (this.filter = props.filter);
  }
}

export class ValidateListUsersInput {
  static validate(input: ListUsersInputConstructorProps) {
    const listUsersInput = new ListUsersInput(input);
    const errors = validateSync(listUsersInput);
    return errors;
  }
}
