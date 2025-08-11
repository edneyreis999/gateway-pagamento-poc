import { IUseCase } from '../../../../shared/application/use-case.interface';
import { EntityValidationError } from '../../../../shared/domain/validators/validation.error';
import { User } from '../../../domain/user.entity';
import { IUserRepository } from '../../../domain/user.repository';
import { UserOutputMapper, type UserOutput } from '../_user-shared/user-output';
import { ValidateCreateUserInput, type CreateUserInput } from './create-user.input';

export class CreateUserUseCase implements IUseCase<CreateUserInput, CreateUserOutput> {
  constructor(private readonly userRepo: IUserRepository) {}

  async execute(input: CreateUserInput): Promise<CreateUserOutput> {
    const errors = ValidateCreateUserInput.validate(input);
    if (errors.length > 0) {
      throw errors;
    }

    const user = User.create(input);

    if (user.notification.hasErrors()) {
      throw new EntityValidationError(user.notification.toJSON());
    }

    await this.userRepo.insert(user);

    return UserOutputMapper.toOutput(user);
  }
}

export type CreateUserOutput = UserOutput;
