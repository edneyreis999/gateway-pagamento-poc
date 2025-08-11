import type { ValidationError } from 'class-validator';
import { NotFoundError } from '../../../../shared/domain/errors/not-found.error';
import { Uuid } from '../../../../shared/domain/value-objects/uuid.vo';
import { User } from '../../../domain/user.entity';
import { UserInMemoryRepository } from '../../../infra/db/in-memory/user-in-memory.repository';
import { DeleteUserUseCase } from './delete-user.use-case';

describe('DeleteUserUseCase Unit Tests', () => {
  let useCase: DeleteUserUseCase;
  let repository: UserInMemoryRepository;

  beforeEach(() => {
    repository = new UserInMemoryRepository();
    useCase = new DeleteUserUseCase(repository);
  });

  it('should throws error when entity not found', async () => {
    const uuid = new Uuid();

    await expect(() => useCase.execute({ id: uuid.id })).rejects.toThrow(
      new NotFoundError(uuid.id, User),
    );
  });

  it('should throw validation error when id is null or undefined', async () => {
    await useCase.execute({ id: null as any }).catch((errors: ValidationError[]) => {
      expect(errors).toBeInstanceOf(Array);
      expect(errors).toHaveLength(1);
      expect(errors[0].constraints).toEqual({
        isUuid: 'id must be a UUID',
      });
    });

    await useCase.execute({ id: undefined as any }).catch((errors: ValidationError[]) => {
      expect(errors).toBeInstanceOf(Array);
      expect(errors).toHaveLength(1);
      expect(errors[0].constraints).toEqual({
        isUuid: 'id must be a UUID',
      });
    });
  });

  it(`should throw validation error when id is not a valid uuid`, async () => {
    await useCase.execute({ id: 'invalid-uuid' }).catch((errors: ValidationError[]) => {
      expect(errors).toBeInstanceOf(Array);
      expect(errors).toHaveLength(1);
      expect(errors[0].constraints).toEqual({
        isUuid: 'id must be a UUID',
      });
    });
  });

  it('should thorw validation error when id is null or undefined', async () => {
    await useCase.execute({ id: null as any }).catch((errors: ValidationError[]) => {
      expect(errors).toBeInstanceOf(Array);
      expect(errors).toHaveLength(1);
      expect(errors[0].constraints).toEqual({
        isUuid: 'id must be a UUID',
      });
    });

    await useCase.execute({ id: undefined as any }).catch((errors: ValidationError[]) => {
      expect(errors).toBeInstanceOf(Array);
      expect(errors).toHaveLength(1);
      expect(errors[0].constraints).toEqual({
        isUuid: 'id must be a UUID',
      });
    });
  });

  it('should delete a user', async () => {
    const items = [new User({ displayName: 'test 1' })];
    repository.items = items;
    await useCase.execute({
      id: items[0].userId.id,
    });
    expect(repository.items).toHaveLength(0);
  });
});
