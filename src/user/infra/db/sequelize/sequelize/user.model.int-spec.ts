import { DataType } from 'sequelize-typescript';
import { setupSequelize } from '../../../../../shared/infra/testing/helpers';
import { UserModel } from './user.model';

describe('UserModel Integration Tests', () => {
  setupSequelize({ models: [UserModel] });

  test('mapping props', () => {
    const attributesMap = UserModel.getAttributes();
    const attributes = Object.keys(UserModel.getAttributes());

    expect(attributes).toStrictEqual([
      'userId',
      'displayName',
      'dustBalance',
      'isActive',
      'createdAt',
    ]);

    const userIdAttr = attributesMap.userId;
    expect(userIdAttr).toMatchObject({
      field: 'userId',
      fieldName: 'userId',
      primaryKey: true,
      type: DataType.UUID(),
    });

    const displayNameAttr = attributesMap.displayName;
    expect(displayNameAttr).toMatchObject({
      field: 'displayName',
      fieldName: 'displayName',
      allowNull: false,
      type: DataType.STRING(255),
    });

    const dustBalanceAttr = attributesMap.dustBalance;
    expect(dustBalanceAttr).toMatchObject({
      field: 'dustBalance',
      fieldName: 'dustBalance',
      allowNull: false,
      type: DataType.NUMBER(),
    });

    const isActiveAttr = attributesMap.isActive;
    expect(isActiveAttr).toMatchObject({
      field: 'isActive',
      fieldName: 'isActive',
      allowNull: false,
      type: DataType.BOOLEAN(),
    });

    const createdAtAttr = attributesMap.createdAt;
    expect(createdAtAttr).toMatchObject({
      field: 'createdAt',
      fieldName: 'createdAt',
      allowNull: false,
      type: DataType.DATE(3),
    });
  });

  test('create', async () => {
    //arrange
    const arrange = {
      userId: '9366b7dc-2d71-4799-b91c-c64adb205104',
      displayName: 'test',
      dustBalance: 100,
      isActive: true,
      createdAt: new Date(),
    };

    //act
    const user = await UserModel.create(arrange);

    //assert
    expect(user.toJSON()).toStrictEqual(arrange);
  });
});
