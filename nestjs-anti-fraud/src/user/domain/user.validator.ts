import { Max, MaxLength, Min } from 'class-validator';
import { ClassValidatorFields } from '../../shared/domain/validators/class-validator-fields';
import { Notification } from '../../shared/domain/validators/notification';
import { User } from './user.entity';

export type ValidFields = 'displayName' | 'dustBalance';
const allValidFields: ValidFields[] = ['displayName', 'dustBalance'];
export class UserRules {
  @MaxLength(30, {
    message: 'Display name must be less than 30 characters',
    groups: ['displayName'],
  })
  declare displayName: string;

  @Min(0, { message: 'Dust balance must be greater than 0', groups: ['dustBalance'] })
  @Max(999999, { message: 'Dust balance must be less than 999999', groups: ['dustBalance'] })
  declare dustBalance: number;

  constructor(entity: User) {
    Object.assign(this, entity);
  }
}

export class UserValidator extends ClassValidatorFields<ValidFields> {
  validate(notification: Notification, data: any, fields?: ValidFields[]): boolean {
    const selectedFields = fields?.length ? fields : allValidFields;
    const userRules = new UserRules(data);
    return super.validate(notification, userRules, selectedFields);
  }
}

export class UserValidatorFactory {
  static create() {
    return new UserValidator();
  }
}
