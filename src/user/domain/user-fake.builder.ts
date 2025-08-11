import { Chance } from 'chance';
import { Uuid } from '../../shared/domain/value-objects/uuid.vo';
import { User } from './user.entity';

type PropOrFactory<T> = T | ((index: number) => T);

export class UserFakeBuilder<TBuild = any> {
  // auto generated in entity
  private _userId: PropOrFactory<Uuid> | undefined = undefined;
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  private _displayName: PropOrFactory<string> = _index => this.chance.name({ prefix: true });
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  private _dustBalance: PropOrFactory<number | null> = _index =>
    this.chance.integer({ min: 0, max: 5000 });
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  private _isActive: PropOrFactory<boolean> = _index => true;
  // auto generated in entity
  private _createdAt: PropOrFactory<Date> | undefined = undefined;

  private countObjs;

  static aUser() {
    return new UserFakeBuilder<User>();
  }

  static theUsers(countObjs: number) {
    return new UserFakeBuilder<User[]>(countObjs);
  }

  private chance: Chance.Chance;

  private constructor(countObjs: number = 1) {
    this.countObjs = countObjs;
    this.chance = Chance();
  }

  withUuid(valueOrFactory: PropOrFactory<Uuid>) {
    this._userId = valueOrFactory;
    return this;
  }

  withDisplayName(valueOrFactory: PropOrFactory<string>) {
    this._displayName = valueOrFactory;
    return this;
  }

  withDustBalance(valueOrFactory: PropOrFactory<number | null>) {
    this._dustBalance = valueOrFactory;
    return this;
  }

  activate() {
    this._isActive = true;
    return this;
  }

  deactivate() {
    this._isActive = false;
    return this;
  }

  withCreatedAt(valueOrFactory: PropOrFactory<Date>) {
    this._createdAt = valueOrFactory;
    return this;
  }

  withInvaliddisplayNameTooLong(value?: string) {
    this._displayName = value ?? this.chance.word({ length: 256 });
    return this;
  }

  build(): TBuild {
    const users = new Array(this.countObjs).fill(undefined).map((_, index) => {
      const user = new User({
        userId: !this._userId ? undefined : this.callFactory(this._userId, index),
        displayName: this.callFactory(this._displayName, index),
        dustBalance: this.callFactory(this._dustBalance, index),
        isActive: this.callFactory(this._isActive, index),
        ...(this._createdAt && {
          createdAt: this.callFactory(this._createdAt, index),
        }),
      });
      user.validate();
      return user;
    });
    return this.countObjs === 1 ? (users[0] as any) : users;
  }

  get userId() {
    return this.getValue('userId');
  }

  get displayName() {
    return this.getValue('displayName');
  }

  get dustBalance() {
    return this.getValue('dustBalance');
  }

  get isActive() {
    return this.getValue('isActive');
  }

  get createdAt() {
    return this.getValue('createdAt');
  }

  private getValue(prop: any) {
    const optional = ['userId', 'createdAt'];
    const privateProp = `_${prop}` as keyof this;
    if (!this[privateProp] && optional.includes(prop)) {
      throw new Error(`Property ${prop} not have a factory, use 'with' methods`);
    }
    return this.callFactory(this[privateProp], 0);
  }

  private callFactory(factoryOrValue: PropOrFactory<any>, index: number) {
    return typeof factoryOrValue === 'function' ? factoryOrValue(index) : factoryOrValue;
  }
}
