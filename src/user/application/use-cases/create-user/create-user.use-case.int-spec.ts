import { Uuid } from '../../../../shared/domain/value-objects/uuid.vo';
import { setupSequelize } from '../../../../shared/infra/testing/helpers';
import { UserSequelizeRepository } from '../../../infra/db/sequelize/sequelize/user-sequelize.repository';
import { UserModel } from '../../../infra/db/sequelize/sequelize/user.model';
import { CreateUserUseCase } from './create-user.use-case';

describe('CreateUserUseCase Integration Tests', () => {
  let useCase: CreateUserUseCase;
  let repository: UserSequelizeRepository;

  setupSequelize({ models: [UserModel] });

  beforeEach(() => {
    repository = new UserSequelizeRepository(UserModel);
    useCase = new CreateUserUseCase(repository);
  });

  describe('Successful User Creation', () => {
    it('should create a user with default properties', async () => {
      const output = await useCase.execute({ displayName: 'test' });
      const entity = await repository.findById(new Uuid(output.id));

      expect(output).toStrictEqual({
        id: entity!.userId.id,
        displayName: 'test',
        dustBalance: 0,
        isActive: true,
        createdAt: entity!.createdAt,
      });

      expect(entity).toBeDefined();
      expect(entity!.displayName).toBe('test');
      expect(entity!.dustBalance).toBe(0);
      expect(entity!.isActive).toBe(true);
    });

    it('should create a user with custom dustBalance', async () => {
      const output = await useCase.execute({ displayName: 'test', dustBalance: 245 });
      const entity = await repository.findById(new Uuid(output.id));

      expect(output).toStrictEqual({
        id: entity!.userId.id,
        displayName: 'test',
        dustBalance: 245,
        isActive: true,
        createdAt: entity!.createdAt,
      });

      expect(entity).toBeDefined();
      expect(entity!.dustBalance).toBe(245);
    });

    it('should create a user with custom isActive true', async () => {
      const output = await useCase.execute({
        displayName: 'test',
        dustBalance: 245,
        isActive: true,
      });
      const entity = await repository.findById(new Uuid(output.id));

      expect(output).toStrictEqual({
        id: entity!.userId.id,
        displayName: 'test',
        dustBalance: 245,
        isActive: true,
        createdAt: entity!.createdAt,
      });

      expect(entity).toBeDefined();
      expect(entity!.isActive).toBe(true);
    });

    it('should create a user with custom isActive false', async () => {
      const output = await useCase.execute({
        displayName: 'test',
        dustBalance: 245,
        isActive: false,
      });
      const entity = await repository.findById(new Uuid(output.id));

      expect(output).toStrictEqual({
        id: entity!.userId.id,
        displayName: 'test',
        dustBalance: 245,
        isActive: false,
        createdAt: entity!.createdAt,
      });

      expect(entity).toBeDefined();
      expect(entity!.isActive).toBe(false);
    });
  });
});
