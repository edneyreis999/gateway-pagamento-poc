import { IRepository } from '../../shared/domain/repository/repository-interface';
import { Uuid } from '../../shared/domain/value-objects/uuid.vo';
import { User } from '../domain/user.entity';

export interface IUserRepository extends IRepository<User, Uuid> {}
