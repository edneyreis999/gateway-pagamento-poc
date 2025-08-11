import { NotFoundError } from '../../../../shared/domain/errors/not-found.error';
import { Uuid } from '../../../../shared/domain/value-objects/uuid.vo';
import { setupSequelize } from '../../../../shared/infra/testing/helpers';
import { User } from '../../../domain/user.entity';
import { UserSequelizeRepository } from '../../../infra/db/sequelize/sequelize/user-sequelize.repository';
import { UserModel } from '../../../infra/db/sequelize/sequelize/user.model';
import { GetUserUseCase } from './get-user.use-case';
describe('GetUserUseCase Integration Tests', () => {
  let useCase: GetUserUseCase;
  let repository: UserSequelizeRepository;

  setupSequelize({ models: [UserModel] });

  beforeEach(() => {
    repository = new UserSequelizeRepository(UserModel);
    useCase = new GetUserUseCase(repository);
  });

  it('should throws error when entity not found', async () => {
    const uuid = new Uuid();
    await expect(() => useCase.execute({ id: uuid.id })).rejects.toThrow(
      new NotFoundError(uuid.id, User),
    );
  });

  it('should returns a user', async () => {
    const user = User.fake().aUser().build();
    await repository.insert(user);

    const spyUpdate = jest.spyOn(repository, 'findById');
    const output = await useCase.execute({ id: user.userId.id });

    expect(spyUpdate).toHaveBeenCalledTimes(1);
    expect(output).toStrictEqual({
      id: user.userId.id,
      displayName: user.displayName,
      dustBalance: user.dustBalance,
      isActive: user.isActive,
      createdAt: user.createdAt,
    });
  });
});
