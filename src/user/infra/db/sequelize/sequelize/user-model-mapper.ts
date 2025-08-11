import { EntityValidationError } from '../../../../../shared/domain/validators/validation.error';
import { Uuid } from '../../../../../shared/domain/value-objects/uuid.vo';
import { User } from '../../../../domain/user.entity';
import { UserModel } from './user.model';

export class UserModelMapper {
  static toModel(entity: User): UserModel {
    return UserModel.build({
      userId: entity.userId.id,
      displayName: entity.displayName,
      dustBalance: entity.dustBalance,
      isActive: entity.isActive,
      createdAt: entity.createdAt,
    });
  }

  static toEntity(model: UserModel): User {
    const entity = new User({
      userId: new Uuid(model.userId),
      displayName: model.displayName,
      dustBalance: model.dustBalance,
      isActive: model.isActive,
      createdAt: model.createdAt,
    });

    entity.validate();
    if (entity.notification.hasErrors()) {
      throw new EntityValidationError(entity.notification.toJSON());
    }
    return entity;
  }
}
