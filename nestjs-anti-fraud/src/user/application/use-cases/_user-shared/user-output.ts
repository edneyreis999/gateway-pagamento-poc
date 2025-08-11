import { User } from '../../../domain/user.entity';

export type UserOutput = {
  id: string;
  displayName: string;
  dustBalance: number;
  isActive: boolean;
  createdAt: Date;
};

export class UserOutputMapper {
  static toOutput(entity: User): UserOutput {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { userId, ...otherProps } = entity.toJSON();
    return {
      id: entity.userId.id,
      ...otherProps,
    };
  }
}
