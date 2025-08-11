import { validateSync } from 'class-validator';
import { Notification } from './notification';
import { IValidatorFields } from './validator-fields-interface';

export abstract class ClassValidatorFields<T extends string> implements IValidatorFields {
  validate(notification: Notification, data: any, fields: T[]): boolean {
    const stringFields = fields.map(field => String(field));
    const errors = validateSync(data, {
      groups: stringFields,
    });
    if (errors.length) {
      for (const error of errors) {
        const field = error.property;
        Object.values(error.constraints!).forEach(message => {
          notification.addError(message, field);
        });
      }
    }
    return !errors.length;
  }
}
