import { NotFoundError } from '../../../../shared/domain/errors/not-found.error';
import { Uuid } from '../../../../shared/domain/value-objects/uuid.vo';
import { setupSequelize } from '../../../../shared/infra/testing/helpers';
import { User } from '../../../domain/user.entity';
import { UserSequelizeRepository } from '../../../infra/db/sequelize/sequelize/user-sequelize.repository';
import { UserModel } from '../../../infra/db/sequelize/sequelize/user.model';
import { DeleteUserUseCase } from './delete-user.use-case';

describe('DeleteUserUseCase Integration Tests', () => {
  let useCase: DeleteUserUseCase;
  let repository: UserSequelizeRepository;

  setupSequelize({ models: [UserModel] });

  beforeEach(() => {
    repository = new UserSequelizeRepository(UserModel);
    useCase = new DeleteUserUseCase(repository);
  });

  it('should throws error when entity not found', async () => {
    const uuid = new Uuid();
    await expect(() => useCase.execute({ id: uuid.id })).rejects.toThrow(
      new NotFoundError(uuid.id, User),
    );
  });

  it('should delete a user', async () => {
    const user = User.fake().aUser().build();
    await repository.insert(user);

    const updateSpy = jest.spyOn(repository, 'delete');
    await useCase.execute({
      id: user.userId.id,
    });
    await expect(repository.findById(user.userId)).resolves.toBeNull();
    expect(updateSpy).toHaveBeenCalledTimes(1);
  });
});
