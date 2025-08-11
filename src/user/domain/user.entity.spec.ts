/* eslint-disable @typescript-eslint/no-explicit-any */
import { Uuid } from '../../shared/domain/value-objects/uuid.vo';
import { User } from './user.entity';

describe('User Unit Tests', () => {
  beforeEach(() => {
    User.prototype.validate = jest.fn().mockImplementation(User.prototype.validate);
  });

  describe('constructor', () => {
    test('should create a user with default values', () => {
      const user = new User({
        displayName: 'John Doe',
      });
      expect(user.userId).toBeInstanceOf(Uuid);
      expect(user.displayName).toBe('John Doe');
      expect(user.dustBalance).toBe(0);
      expect(user.isActive).toBeTruthy();
      expect(user.createdAt).toBeInstanceOf(Date);
    });

    test('should create a user with all values', () => {
      const createdAt = new Date();
      const user = new User({
        displayName: 'John Doe',
        dustBalance: 100,
        isActive: false,
        createdAt,
      });
      expect(user.userId).toBeInstanceOf(Uuid);
      expect(user.displayName).toBe('John Doe');
      expect(user.dustBalance).toBe(100);
      expect(user.isActive).toBeFalsy();
      expect(user.createdAt).toBe(createdAt);
    });
  });

  describe('create command', () => {
    test('should create a user', () => {
      const user = User.create({
        displayName: 'John Doe',
      });
      expect(user.userId).toBeInstanceOf(Uuid);
      expect(user.displayName).toBe('John Doe');
      expect(user.dustBalance).toBe(0);
      expect(user.isActive).toBe(true);
      expect(user.createdAt).toBeInstanceOf(Date);
      expect(User.prototype.validate).toHaveBeenCalledTimes(1);
      console.log(user.notification.toJSON());
      expect(user.notification.hasErrors()).toBe(false);
    });

    test('should create a user with all properties', () => {
      const user = User.create({
        displayName: 'John Doe',
        dustBalance: 100,
        isActive: false,
      });
      expect(user.userId).toBeInstanceOf(Uuid);
      expect(user.displayName).toBe('John Doe');
      expect(user.dustBalance).toBe(100);
      expect(user.isActive).toBe(false);
      expect(user.createdAt).toBeInstanceOf(Date);
      expect(User.prototype.validate).toHaveBeenCalledTimes(1);
      expect(user.notification.hasErrors()).toBe(false);
    });
  });

  describe('userId field', () => {
    const arrange = [{ userId: null }, { userId: undefined }, { userId: new Uuid() }];
    test.each(arrange)('userId = %j', ({ userId }) => {
      const user = new User({
        displayName: 'John Doe',
        userId: userId as any,
      });
      expect(user.userId).toBeInstanceOf(Uuid);
      if (userId instanceof Uuid) {
        expect(user.userId).toBe(userId);
      }
    });
  });

  test('should change displayName', () => {
    const user = User.create({
      displayName: 'Joao',
    });
    user.changeDisplayName('Jane Doe');
    expect(user.displayName).toBe('Jane Doe');
    expect(User.prototype.validate).toHaveBeenCalledTimes(2);
    expect(user.notification.hasErrors()).toBe(false);
  });

  test('should add dust to dustBalance', () => {
    const user = User.create({
      displayName: 'John Doe',
    });
    user.addDust(50);
    expect(user.dustBalance).toBe(50);
    expect(User.prototype.validate).toHaveBeenCalledTimes(2);
    expect(user.notification.hasErrors()).toBe(false);
  });

  test('should subtract dust from dustBalance', () => {
    const user = User.create({
      displayName: 'John Doe',
      dustBalance: 100,
    });
    user.subtractDust(50);
    expect(user.dustBalance).toBe(50);
    expect(User.prototype.validate).toHaveBeenCalledTimes(2);
    expect(user.notification.hasErrors()).toBe(false);
  });

  test('should throw error when subtracting negative dust', () => {
    const user = User.create({
      displayName: 'John Doe',
      dustBalance: 50,
    });
    user.subtractDust(60);
    expect(User.prototype.validate).toHaveBeenCalledTimes(2);
    expect(user.notification.hasErrors()).toBe(true);
    expect(Object.fromEntries(user.notification.errors)).toEqual({
      dustBalance: ['Dust balance must be greater than 0'],
    });
  });

  test('should throw error when dust is more then 999999', () => {
    const user = User.create({
      displayName: 'John Doe',
      dustBalance: 50,
    });
    user.addDust(999999);
    expect(user.notification.hasErrors()).toBe(true);
    expect(Object.fromEntries(user.notification.errors)).toEqual({
      dustBalance: ['Dust balance must be less than 999999'],
    });
  });

  test('should throw error when subtracting more dust than available', () => {
    const user = User.create({
      displayName: 'John Doe',
      dustBalance: 50,
    });
    user.subtractDust(100);
    expect(user.notification.hasErrors()).toBe(true);
    expect(Object.fromEntries(user.notification.errors)).toEqual({
      dustBalance: ['Dust balance must be greater than 0'],
    });
  });

  test('should activate a user', () => {
    const user = User.create({
      displayName: 'John Doe',
      isActive: false,
    });
    user.activate();
    expect(user.isActive).toBe(true);
  });

  test('should deactivate a user', () => {
    const user = User.create({
      displayName: 'John Doe',
      isActive: true,
    });
    user.deactivate();
    expect(user.isActive).toBe(false);
  });

  it('should get user by entityId', () => {
    const user = User.create({
      displayName: 'John Doe',
    });
    expect(user.entityId).toBe(user.userId);
  });

  it('should transform user to JSON', () => {
    const user = User.create({
      displayName: 'John Doe',
    });
    expect(user.toJSON()).toEqual({
      userId: user.userId.id,
      displayName: 'John Doe',
      dustBalance: 0,
      isActive: true,
      createdAt: user.createdAt,
    });
  });
});

describe('User Validator', () => {
  describe('create command', () => {
    test('should throw an error for invalid displayName', () => {
      const user = User.create({ displayName: 't'.repeat(31) });
      expect(user.notification.hasErrors()).toBe(true);
      expect(user.notification).notificationContainsErrorMessages([
        {
          displayName: ['Display name must be less than 30 characters'],
        },
      ]);
    });

    test('should throw an error for dustBalance lass then 0', () => {
      const user = User.create({ displayName: 'John Doe', dustBalance: -1 });
      expect(user.notification.hasErrors()).toBe(true);
      expect(user.notification).notificationContainsErrorMessages([
        { dustBalance: ['Dust balance must be greater than 0'] },
      ]);
    });

    test('should throw an error for dustBalance bigger then 999999', () => {
      const user = User.create({ displayName: 'John Doe', dustBalance: 1000000 });
      expect(user.notification.hasErrors()).toBe(true);
      expect(user.notification).notificationContainsErrorMessages([
        { dustBalance: ['Dust balance must be less than 999999'] },
      ]);
    });
  });

  describe('changeDisplayName method', () => {
    test('should throw an error for invalid displayName', () => {
      const user = User.create({ displayName: 'John Doe' });
      user.changeDisplayName('t'.repeat(256));

      expect(user.notification.hasErrors()).toBe(true);
      expect(user.notification).notificationContainsErrorMessages([
        {
          displayName: ['Display name must be less than 30 characters'],
        },
      ]);
    });

    test('should not throw an error for valid displayName', () => {
      const user = User.create({ displayName: 'John Doe' });
      user.changeDisplayName('Jane Doe');

      expect(user.notification.hasErrors()).toBe(false);
    });

    test('should not throw an error for the same displayName', () => {
      const user = User.create({ displayName: 'John Doe' });
      user.changeDisplayName('John Doe');

      expect(user.notification.hasErrors()).toBe(false);
    });
    // test dust balance

    test('should throw an error for add dust more then 999999', () => {
      const user = User.create({ displayName: 'John Doe' });
      user.addDust(1000000);

      expect(user.notification.hasErrors()).toBe(true);
      expect(user.notification).notificationContainsErrorMessages([
        { dustBalance: ['Dust balance must be less than 999999'] },
      ]);
    });

    test('should throw an error for subtract dust lass then 0', () => {
      const user = User.create({ displayName: 'John Doe' });
      user.addDust(500);
      user.subtractDust(600);

      expect(user.notification.hasErrors()).toBe(true);
      expect(user.notification).notificationContainsErrorMessages([
        { dustBalance: ['Dust balance must be greater than 0'] },
      ]);
    });
  });
});
