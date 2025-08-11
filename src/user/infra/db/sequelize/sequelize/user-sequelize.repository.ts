import { Op } from 'sequelize';
import { NotFoundError } from '../../../../../shared/domain/errors/not-found.error';
import type { Uuid } from '../../../../../shared/domain/value-objects/uuid.vo';
import { User } from '../../../../domain/user.entity';
import {
  type IUserRepository,
  type UserSearchParams,
  UserSearchResult,
} from '../../../../domain/user.repository';
import { UserModelMapper } from './user-model-mapper';
import { UserModel } from './user.model';

export class UserSequelizeRepository implements IUserRepository {
  sortableFields: string[] = ['displayName', 'createdAt'];

  constructor(private userModel: typeof UserModel) {}

  async insert(entity: User): Promise<void> {
    const modelProps = UserModelMapper.toModel(entity);
    await this.userModel.create(modelProps.toJSON());
  }

  async bulkInsert(entities: User[]): Promise<void> {
    const modelsProps = entities.map(entity => UserModelMapper.toModel(entity).toJSON());
    await this.userModel.bulkCreate(modelsProps);
  }

  async update(entity: User): Promise<void> {
    const id = entity.userId.id;
    const model = await this._get(id);
    if (!model) {
      throw new NotFoundError(id, this.getEntity());
    }
    const modelProps = UserModelMapper.toModel(entity);
    await this.userModel.update(modelProps.toJSON(), {
      where: { userId: id },
    });
  }

  async delete(userId: Uuid): Promise<void> {
    const id = userId.id;
    const model = await this._get(id);
    if (!model) {
      throw new NotFoundError(id, this.getEntity());
    }
    await this.userModel.destroy({ where: { userId: id } });
  }

  async findById(entity_id: Uuid): Promise<User | null> {
    const model = await this._get(entity_id.id);

    return model ? UserModelMapper.toEntity(model) : null;
  }

  private async _get(id: string) {
    return await this.userModel.findByPk(id);
  }

  async findAll(): Promise<User[]> {
    const models = await this.userModel.findAll();
    return models.map(model => {
      return UserModelMapper.toEntity(model);
    });
  }

  async search(props: UserSearchParams): Promise<UserSearchResult> {
    const offset = (props.page - 1) * props.per_page;
    const limit = props.per_page;
    const sortDirection = props.sort_dir ?? 'DESC';

    const { rows: models, count } = await this.userModel.findAndCountAll({
      ...(props.filter && {
        where: {
          displayName: { [Op.like]: `%${props.filter}%` },
        },
      }),
      ...(props.sort && this.sortableFields.includes(props.sort)
        ? { order: [[props.sort, sortDirection]] }
        : { order: [['createdAt', 'desc']] }),
      offset,
      limit,
    });
    return new UserSearchResult({
      items: models.map(model => {
        return UserModelMapper.toEntity(model);
      }),
      current_page: props.page,
      per_page: props.per_page,
      total: count,
    });
  }

  getEntity(): new (...args: any[]) => User {
    return User;
  }
}
