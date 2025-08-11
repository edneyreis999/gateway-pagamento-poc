import type { ValidationError } from 'class-validator';
import { NotFoundError } from '../../../../shared/domain/errors/not-found.error';
import { EntityValidationError } from '../../../../shared/domain/validators/validation.error';
import { Uuid } from '../../../../shared/domain/value-objects/uuid.vo';
import { User } from '../../../domain/user.entity';
import { UserInMemoryRepository } from '../../../infra/db/in-memory/user-in-memory.repository';
import { AddDustUserUseCase } from './dust-add-user.use-case';

describe('AddDustUserUseCase Unit Tests', () => {
  let useCase: AddDustUserUseCase;
  let repository: UserInMemoryRepository;

  beforeEach(() => {
    repository = new UserInMemoryRepository();
    useCase = new AddDustUserUseCase(repository);
  });

  it('should throws error when entity not found', async () => {
    const uuid = new Uuid();

    await expect(() => useCase.execute({ id: uuid.id, dust: 800 })).rejects.toThrow(
      new NotFoundError(uuid.id, User),
    );
  });

  it('should add dust to a user', async () => {
    const spyAddDust = jest.spyOn(repository, 'update');
    const initialDust = 500;
    let dustBalance = initialDust;
    const entity = User.fake().aUser().withDisplayName('test').withDustBalance(initialDust).build();
    repository.items = [entity];

    let output = await useCase.execute({
      id: entity.userId.id,
      dust: 800,
    });
    expect(spyAddDust).toHaveBeenCalledTimes(1);
    dustBalance += 800;
    expect(output).toStrictEqual({
      id: entity.userId.id,
      displayName: 'test',
      dustBalance: dustBalance,
      isActive: true,
      createdAt: entity.createdAt,
    });

    type Arrange = {
      input: {
        id: string;
        dust: number;
      };
      expected: {
        id: string;
        displayName: string;
        dustBalance: number | null;
        isActive: boolean;
        createdAt: Date;
      };
    };
    const arrange: Arrange[] = [
      {
        input: {
          id: entity.userId.id,
          dust: 0,
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          dustBalance: dustBalance,
          isActive: true,
          createdAt: entity.createdAt,
        },
      },
      {
        input: {
          id: entity.userId.id,
          dust: 300,
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          dustBalance: dustBalance + 300,
          isActive: true,
          createdAt: entity.createdAt,
        },
      },
    ];

    for (const i of arrange) {
      output = await useCase.execute({
        id: i.input.id,
        dust: i.input.dust,
      });
      expect(output).toStrictEqual({
        id: entity.userId.id,
        displayName: entity.displayName,
        dustBalance: i.expected.dustBalance,
        isActive: entity.isActive,
        createdAt: entity.createdAt,
      });
    }
  });

  it(`should't add dust to a user when dustBalance is bigger then 999999`, async () => {
    const entity = User.fake().aUser().withDisplayName('test').withDustBalance(500).build();
    repository.items = [entity];
    await useCase
      .execute({ id: entity.userId.id, dust: 1000000 })
      .catch((error: EntityValidationError) => {
        expect(error).toBeInstanceOf(EntityValidationError);
        expect(error.error).toEqual([{ dustBalance: ['Dust balance must be less than 999999'] }]);
      });
  });

  it(`should throw validation error when dust is less then 0`, async () => {
    const entity = User.fake().aUser().withDisplayName('test').withDustBalance(500).build();
    expect.assertions(3);
    await repository.insert(entity);
    await useCase
      .execute({ id: entity.userId.id, dust: -10 })
      .catch((errors: ValidationError[]) => {
        expect(errors).toBeInstanceOf(Array);
        expect(errors).toHaveLength(1);
        expect(errors[0].constraints).toEqual({
          min: 'dust must not be less than 0',
        });
      });
  });

  it(`should throw validation error when dust is not a number`, async () => {
    await useCase
      .execute({ id: new Uuid().id, dust: 'invalid-number' as any })
      .catch((errors: ValidationError[]) => {
        expect(errors).toBeInstanceOf(Array);
        expect(errors).toHaveLength(1);
        expect(errors[0].constraints).toEqual({
          isNumber: 'dust must be a number conforming to the specified constraints',
          min: 'dust must not be less than 0',
        });
      });
  });

  it('should throw validation error when id is null or undefined', async () => {
    await useCase.execute({ id: null as any, dust: 500 }).catch((errors: ValidationError[]) => {
      expect(errors).toBeInstanceOf(Array);
      expect(errors).toHaveLength(1);
      expect(errors[0].constraints).toEqual({
        isUuid: 'id must be a UUID',
      });
    });

    await useCase
      .execute({ id: undefined as any, dust: 500 })
      .catch((errors: ValidationError[]) => {
        expect(errors).toBeInstanceOf(Array);
        expect(errors).toHaveLength(1);
        expect(errors[0].constraints).toEqual({
          isUuid: 'id must be a UUID',
        });
      });
  });

  it(`should throw validation error when id is not a valid uuid`, async () => {
    await useCase.execute({ id: 'invalid-uuid', dust: 500 }).catch((errors: ValidationError[]) => {
      expect(errors).toBeInstanceOf(Array);
      expect(errors).toHaveLength(1);
      expect(errors[0].constraints).toEqual({
        isUuid: 'id must be a UUID',
      });
    });
  });

  it('should thorw validation error when id is null or undefined', async () => {
    await useCase.execute({ id: null as any, dust: 500 }).catch((errors: ValidationError[]) => {
      expect(errors).toBeInstanceOf(Array);
      expect(errors).toHaveLength(1);
      expect(errors[0].constraints).toEqual({
        isUuid: 'id must be a UUID',
      });
    });

    await useCase
      .execute({ id: undefined as any, dust: 500 })
      .catch((errors: ValidationError[]) => {
        expect(errors).toBeInstanceOf(Array);
        expect(errors).toHaveLength(1);
        expect(errors[0].constraints).toEqual({
          isUuid: 'id must be a UUID',
        });
      });
  });
});
