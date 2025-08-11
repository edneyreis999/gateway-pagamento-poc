import { Chance } from 'chance';
import { Uuid } from '../../shared/domain/value-objects/uuid.vo';
import { UserFakeBuilder } from './user-fake.builder';

describe('UserFakerBuilder Unit Tests', () => {
  describe('userId prop', () => {
    const faker = UserFakeBuilder.aUser();

    test('should throw error when any with methods has called', () => {
      expect(() => faker.userId).toThrow(
        new Error("Property userId not have a factory, use 'with' methods"),
      );
    });

    test('should be undefined', () => {
      expect(faker['_userId']).toBeUndefined();
    });

    test('withUuid', () => {
      const userId = new Uuid();
      const $this = faker.withUuid(userId);
      expect($this).toBeInstanceOf(UserFakeBuilder);
      expect(faker['_userId']).toBe(userId);

      faker.withUuid(() => userId);
      //@ts-expect-error _userId is a callable
      expect(faker['_userId']()).toBe(userId);

      expect(faker.userId).toBe(userId);
    });

    test('should pass index to userId factory', () => {
      let mockFactory = jest.fn(() => new Uuid());
      faker.withUuid(mockFactory);
      faker.build();
      expect(mockFactory).toHaveBeenCalledTimes(1);

      const userId = new Uuid();
      mockFactory = jest.fn(() => userId);
      const fakerMany = UserFakeBuilder.theUsers(2);
      fakerMany.withUuid(mockFactory);
      fakerMany.build();

      expect(mockFactory).toHaveBeenCalledTimes(2);
      expect(fakerMany.build()[0].userId).toBe(userId);
      expect(fakerMany.build()[1].userId).toBe(userId);
    });
  });

  describe('displayName prop', () => {
    const faker = UserFakeBuilder.aUser();
    test('should be a function', () => {
      expect(typeof faker['_displayName']).toBe('function');
    });

    test('should call the displayName method', () => {
      const chance = Chance();
      const spyWordMethod = jest.spyOn(chance, 'name');
      faker['chance'] = chance;
      faker.build();

      expect(spyWordMethod).toHaveBeenCalled();
    });

    test('withDisplayName', () => {
      const $this = faker.withDisplayName('test displayName');
      expect($this).toBeInstanceOf(UserFakeBuilder);
      expect(faker['_displayName']).toBe('test displayName');

      faker.withDisplayName(() => 'test displayName');
      //@ts-expect-error displayName is callable
      expect(faker['_displayName']()).toBe('test displayName');

      expect(faker.displayName).toBe('test displayName');
    });

    test('should pass index to displayName factory', () => {
      faker.withDisplayName(index => `test displayName ${index}`);
      const user = faker.build();
      expect(user.displayName).toBe(`test displayName 0`);

      const fakerMany = UserFakeBuilder.theUsers(2);
      fakerMany.withDisplayName(index => `test displayName ${index}`);
      const users = fakerMany.build();

      expect(users[0].displayName).toBe(`test displayName 0`);
      expect(users[1].displayName).toBe(`test displayName 1`);
    });

    test('invalid too long case', () => {
      const $this = faker.withInvaliddisplayNameTooLong();
      expect($this).toBeInstanceOf(UserFakeBuilder);
      expect(faker['_displayName'].length).toBe(256);

      const tooLong = 'a'.repeat(256);
      faker.withInvaliddisplayNameTooLong(tooLong);
      expect(faker['_displayName'].length).toBe(256);
      expect(faker['_displayName']).toBe(tooLong);
    });
  });

  describe('dustBalance prop', () => {
    const faker = UserFakeBuilder.aUser();
    test('should be a function', () => {
      expect(typeof faker['_dustBalance']).toBe('function');
    });

    test('should call the paragraph method', () => {
      const chance = Chance();
      const spyParagraphMethod = jest.spyOn(chance, 'integer');
      faker['chance'] = chance;
      faker.build();
      expect(spyParagraphMethod).toHaveBeenCalled();
    });

    test('withDustBalance', () => {
      const $this = faker.withDustBalance(4250);
      expect($this).toBeInstanceOf(UserFakeBuilder);
      expect(faker['_dustBalance']).toBe(4250);

      faker.withDustBalance(() => 4250);
      //@ts-expect-error dustBalance is callable
      expect(faker['_dustBalance']()).toBe(4250);

      expect(faker.dustBalance).toBe(4250);
    });

    test('should pass index to dustBalance factory', () => {
      faker.withDustBalance(index => index);
      const user = faker.build();
      expect(user.dustBalance).toBe(0);

      const fakerMany = UserFakeBuilder.theUsers(2);
      fakerMany.withDustBalance(index => index);
      const users = fakerMany.build();

      expect(users[0].dustBalance).toBe(0);
      expect(users[1].dustBalance).toBe(1);
    });
  });

  describe('isActive prop', () => {
    const faker = UserFakeBuilder.aUser();
    test('should be a function', () => {
      expect(typeof faker['_isActive']).toBe('function');
    });

    test('activate', () => {
      const $this = faker.activate();
      expect($this).toBeInstanceOf(UserFakeBuilder);
      expect(faker['_isActive']).toBe(true);
      expect(faker.isActive).toBe(true);
    });

    test('deactivate', () => {
      const $this = faker.deactivate();
      expect($this).toBeInstanceOf(UserFakeBuilder);
      expect(faker['_isActive']).toBe(false);
      expect(faker.isActive).toBe(false);
    });
  });

  describe('createdAt prop', () => {
    const faker = UserFakeBuilder.aUser();

    test('should throw error when any with methods has called', () => {
      const fakerUser = UserFakeBuilder.aUser();
      expect(() => fakerUser.createdAt).toThrow(
        new Error("Property createdAt not have a factory, use 'with' methods"),
      );
    });

    test('should be undefined', () => {
      expect(faker['_createdAt']).toBeUndefined();
    });

    test('withCreatedAt', () => {
      const date = new Date();
      const $this = faker.withCreatedAt(date);
      expect($this).toBeInstanceOf(UserFakeBuilder);
      expect(faker['_createdAt']).toBe(date);

      faker.withCreatedAt(() => date);
      //@ts-expect-error _createdAt is a callable
      expect(faker['_createdAt']()).toBe(date);
      expect(faker.createdAt).toBe(date);
    });

    test('should pass index to createdAt factory', () => {
      const date = new Date();
      faker.withCreatedAt(index => new Date(date.getTime() + index + 2));
      const user = faker.build();
      expect(user.createdAt.getTime()).toBe(date.getTime() + 2);

      const fakerMany = UserFakeBuilder.theUsers(2);
      fakerMany.withCreatedAt(index => new Date(date.getTime() + index + 2));
      const users = fakerMany.build();

      expect(users[0].createdAt.getTime()).toBe(date.getTime() + 2);
      expect(users[1].createdAt.getTime()).toBe(date.getTime() + 3);
    });
  });

  test('should create a user', () => {
    const faker = UserFakeBuilder.aUser();
    let user = faker.build();

    expect(user.userId).toBeInstanceOf(Uuid);
    expect(typeof user.displayName === 'string').toBeTruthy();
    expect(typeof user.dustBalance === 'number').toBeTruthy();
    expect(user.isActive).toBe(true);
    expect(user.createdAt).toBeInstanceOf(Date);

    const createdAt = new Date();
    const userId = new Uuid();
    user = faker
      .withUuid(userId)
      .withDisplayName('displayName test')
      .withDustBalance(2000)
      .deactivate()
      .withCreatedAt(createdAt)
      .build();

    expect(user.userId.id).toBe(userId.id);
    expect(user.displayName).toBe('displayName test');
    expect(user.dustBalance).toBe(2000);
    expect(user.isActive).toBe(false);
    expect(user.createdAt).toBe(createdAt);
  });

  test('should create many users', () => {
    const faker = UserFakeBuilder.theUsers(2);
    let users = faker.build();

    users.forEach(user => {
      expect(user.userId).toBeInstanceOf(Uuid);
      expect(typeof user.displayName === 'string').toBeTruthy();
      expect(typeof user.dustBalance === 'number').toBeTruthy();
      expect(user.isActive).toBe(true);
      expect(user.createdAt).toBeInstanceOf(Date);
    });

    const createdAt = new Date();
    const userId = new Uuid();
    users = faker
      .withUuid(userId)
      .withDisplayName('displayName test')
      .withDustBalance(2000)
      .deactivate()
      .withCreatedAt(createdAt)
      .build();

    users.forEach(user => {
      expect(user.userId.id).toBe(userId.id);
      expect(user.displayName).toBe('displayName test');
      expect(user.dustBalance).toBe(2000);
      expect(user.isActive).toBe(false);
      expect(user.createdAt).toBe(createdAt);
    });
  });
});
