import type { ValidationError } from 'class-validator';
import { UserInMemoryRepository } from '../../../infra/db/in-memory/user-in-memory.repository';
import { CreateUserUseCase } from './create-user.use-case';

describe('CreateUserUseCase Unit Tests', () => {
  let useCase: CreateUserUseCase;
  let repository: UserInMemoryRepository;

  beforeEach(() => {
    repository = new UserInMemoryRepository();
    useCase = new CreateUserUseCase(repository);
  });

  describe('Successful User Creation', () => {
    it('should create a user with default properties', async () => {
      const spyInsert = jest.spyOn(repository, 'insert');
      const output = await useCase.execute({ displayName: 'test' });

      expect(spyInsert).toHaveBeenCalledTimes(1);
      expect(output).toStrictEqual({
        id: repository.items[0].userId.id,
        displayName: 'test',
        dustBalance: 0,
        isActive: true,
        createdAt: repository.items[0].createdAt,
      });

      const createdUser = repository.items[0];
      expect(createdUser.displayName).toBe('test');
      expect(createdUser.dustBalance).toBe(0);
      expect(createdUser.isActive).toBe(true);
    });

    it('should create a user with custom properties', async () => {
      const spyInsert = jest.spyOn(repository, 'insert');
      const output = await useCase.execute({
        displayName: 'test',
        dustBalance: 230,
        isActive: false,
      });

      expect(spyInsert).toHaveBeenCalledTimes(1);
      expect(output).toStrictEqual({
        id: repository.items[0].userId.id,
        displayName: 'test',
        dustBalance: 230,
        isActive: false,
        createdAt: repository.items[0].createdAt,
      });

      const createdUser = repository.items[0];
      expect(createdUser.displayName).toBe('test');
      expect(createdUser.dustBalance).toBe(230);
      expect(createdUser.isActive).toBe(false);
    });
  });

  describe('Validation Errors', () => {
    it('should throw validation errors for empty displayName', async () => {
      const invalidInput = { displayName: '' };

      await useCase.execute(invalidInput).catch((errors: ValidationError[]) => {
        expect(errors).toBeInstanceOf(Array);
        expect(errors).toHaveLength(1);
        expect(errors[0].constraints).toEqual({
          isNotEmpty: 'displayName should not be empty',
        });
      });
    });

    it('should throw validation errors for negative dustBalance', async () => {
      const invalidInput = { displayName: 'test', dustBalance: -10 };

      await useCase.execute(invalidInput).catch((errors: ValidationError[]) => {
        expect(errors).toBeInstanceOf(Array);
        expect(errors).toHaveLength(1);
        expect(errors[0].constraints).toEqual({
          min: 'dustBalance must not be less than 0',
        });
      });
    });

    it('should throw validation errors for dustBalance exceeding max limit', async () => {
      const invalidInput = { displayName: 'test', dustBalance: 1000000 };

      await useCase.execute(invalidInput).catch((errors: ValidationError[]) => {
        expect(errors).toBeInstanceOf(Array);
        expect(errors).toHaveLength(1);
        expect(errors[0].constraints).toEqual({
          max: 'dustBalance must not be greater than 999999',
        });
      });
    });
  });

  describe('Error Handling', () => {
    it('should handle errors when repository insertion fails', async () => {
      const spyInsert = jest.spyOn(repository, 'insert').mockImplementation(() => {
        throw new Error('Repository error');
      });

      const validInput = { displayName: 'test' };
      await expect(useCase.execute(validInput)).rejects.toThrow('Repository error');

      expect(spyInsert).toHaveBeenCalledTimes(1);
    });
  });
});
