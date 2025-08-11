import { EntityValidationError } from '../../../../../shared/domain/validators/validation.error';
import { Uuid } from '../../../../../shared/domain/value-objects/uuid.vo';
import { setupSequelize } from '../../../../../shared/infra/testing/helpers';

import { User } from '../../../../domain/user.entity';
import { UserModelMapper } from './user-model-mapper';
import { UserModel } from './user.model';

describe('UserModelMapper Integration Tests', () => {
  setupSequelize({ models: [UserModel] });

  it('should throws error when user is invalid', () => {
    expect.assertions(2);
    const model = UserModel.build({
      userId: '9366b7dc-2d71-4799-b91c-c64adb205104',
      displayName: 'a'.repeat(256),
    } as any);
    try {
      UserModelMapper.toEntity(model);
      fail('The user is valid, but it needs throws a EntityValidationError');
    } catch (e) {
      expect(e).toBeInstanceOf(EntityValidationError);
      expect((e as EntityValidationError).error).toMatchObject([
        { displayName: ['Display name must be less than 30 characters'] },
      ]);
    }
  });

  it('should convert a user model to a user aggregate', () => {
    const createdAt = new Date();
    const model = UserModel.build({
      userId: '5490020a-e866-4229-9adc-aa44b83234c4',
      displayName: 'some value',
      dustBalance: 0,
      isActive: true,
      createdAt,
    });
    const aggregate = UserModelMapper.toEntity(model);
    expect(aggregate.toJSON()).toStrictEqual(
      new User({
        userId: new Uuid('5490020a-e866-4229-9adc-aa44b83234c4'),
        displayName: 'some value',
        dustBalance: 0,
        isActive: true,
        createdAt,
      }).toJSON(),
    );
  });

  test('should convert a user aggregate to a user model', () => {
    const createdAt = new Date();
    const aggregate = new User({
      userId: new Uuid('5490020a-e866-4229-9adc-aa44b83234c4'),
      displayName: 'some value',
      dustBalance: 0,
      isActive: true,
      createdAt,
    });
    const model = UserModelMapper.toModel(aggregate);
    expect(model.toJSON()).toStrictEqual({
      userId: '5490020a-e866-4229-9adc-aa44b83234c4',
      displayName: 'some value',
      dustBalance: 0,
      isActive: true,
      createdAt,
    });
  });
});
