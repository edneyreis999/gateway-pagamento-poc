import { Entity } from '../../shared/domain/entity';
import type { ValueObject } from '../../shared/domain/value-object';
import { Uuid } from '../../shared/domain/value-objects/uuid.vo';
import { UserFakeBuilder } from './user-fake.builder';
import { UserValidatorFactory, type ValidFields } from './user.validator';

export type UserConstructorProps = {
  userId?: Uuid;
  displayName: string;
  dustBalance?: number | null;
  isActive?: boolean;
  createdAt?: Date;
};

export type UserCreateCommand = {
  displayName: string;
  dustBalance?: number;
  isActive?: boolean;
};

export class User extends Entity {
  userId: Uuid;
  displayName: string;
  dustBalance: number;
  isActive: boolean;
  createdAt: Date;

  constructor(props: UserConstructorProps) {
    super();
    this.userId = props.userId ?? new Uuid();
    this.displayName = props.displayName;
    this.dustBalance = props.dustBalance ?? 0;
    this.isActive = props.isActive ?? true;
    this.createdAt = props.createdAt ?? new Date();
  }

  get entityId(): ValueObject {
    return this.userId;
  }

  static create(props: UserCreateCommand): User {
    const user = new User(props);
    user.validate(['displayName', 'dustBalance']);
    return user;
  }

  changeDisplayName(displayName: string): void {
    this.displayName = displayName;
    this.validate(['displayName']);
  }

  addDust(amount: number): void {
    this.dustBalance += amount;
    this.validate(['dustBalance']);
  }

  subtractDust(amount: number): void {
    this.dustBalance -= amount;
    this.validate(['dustBalance']);
  }

  activate() {
    this.isActive = true;
  }

  deactivate() {
    this.isActive = false;
  }

  validate(fields?: ValidFields[]) {
    const validator = UserValidatorFactory.create();
    return validator.validate(this.notification, this, fields);
  }

  static fake() {
    return UserFakeBuilder;
  }

  toJSON() {
    return {
      userId: this.userId.id,
      displayName: this.displayName,
      dustBalance: this.dustBalance,
      isActive: this.isActive,
      createdAt: this.createdAt,
    };
  }
}
