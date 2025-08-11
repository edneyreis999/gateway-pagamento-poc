import { IUseCase } from '../../../../shared/application/use-case.interface';
import { Uuid } from '../../../../shared/domain/value-objects/uuid.vo';
import { IUserRepository } from '../../../domain/user.repository';
import { ValidateDeleteUserInput, type DeleteUserInput } from './delete-user.input';

export class DeleteUserUseCase implements IUseCase<DeleteUserInput, void> {
  constructor(private userRepo: IUserRepository) {}

  async execute(input: DeleteUserInput): Promise<void> {
    const errors = ValidateDeleteUserInput.validate(input);
    if (errors.length > 0) {
      throw errors;
    }
    const uuid = new Uuid(input.id);
    await this.userRepo.delete(uuid);
  }
}
