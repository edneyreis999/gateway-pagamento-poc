import type { SortDirection } from '../../../../shared/domain/repository/search-params';
import type { Uuid } from '../../../../shared/domain/value-objects/uuid.vo';
import { InMemorySearchableRepository } from '../../../../shared/infra/db/in-memory/in-memory.repository';
import { User } from '../../../domain/user.entity';

export class UserInMemoryRepository extends InMemorySearchableRepository<User, Uuid> {
  sortableFields: string[] = ['displayName', 'createdAt'];

  getEntity(): new (...args: any[]) => User {
    return User;
  }

  protected async applyFilter(items: User[], filter: string): Promise<User[]> {
    if (!filter) {
      return items;
    }

    return items.filter(i => {
      return i.displayName.toLowerCase().includes(filter.toLowerCase());
    });
  }

  protected applySort(items: User[], sort: string | null, sort_dir: SortDirection | null) {
    return sort
      ? super.applySort(items, sort, sort_dir)
      : super.applySort(items, 'createdAt', 'desc');
  }
}
