import { NotFoundError } from '../../../../shared/domain/errors/not-found.error';
import { Uuid } from '../../../../shared/domain/value-objects/uuid.vo';
import { setupSequelize } from '../../../../shared/infra/testing/helpers';
import { User } from '../../../domain/user.entity';
import { UserSequelizeRepository } from '../../../infra/db/sequelize/sequelize/user-sequelize.repository';
import { UserModel } from '../../../infra/db/sequelize/sequelize/user.model';
import { UpdateUserUseCase } from './update-user.use-case';

describe('UpdateUserUseCase Integration Tests', () => {
  let useCase: UpdateUserUseCase;
  let repository: UserSequelizeRepository;

  setupSequelize({ models: [UserModel] });

  beforeEach(() => {
    repository = new UserSequelizeRepository(UserModel);
    useCase = new UpdateUserUseCase(repository);
  });

  it('should throws error when entity not found', async () => {
    const uuid = new Uuid();
    await expect(() => useCase.execute({ id: uuid.id, displayName: 'fake' })).rejects.toThrow(
      new NotFoundError(uuid.id, User),
    );
  });

  it('should update a user', async () => {
    const entity = User.fake().aUser().build();
    repository.insert(entity);

    const spyUpdate = jest.spyOn(repository, 'update');

    let output = await useCase.execute({
      id: entity.userId.id,
      displayName: 'test',
    });
    expect(spyUpdate).toHaveBeenCalledTimes(1);
    expect(output).toStrictEqual({
      id: entity.userId.id,
      displayName: 'test',
      dustBalance: entity.dustBalance,
      isActive: true,
      createdAt: entity.createdAt,
    });

    type Arrange = {
      input: {
        id: string;
        displayName: string;
        isActive?: boolean;
      };
      expected: {
        id: string;
        displayName: string;
        isActive: boolean;
        createdAt: Date;
      };
    };
    const arrange: Arrange[] = [
      {
        input: {
          id: entity.userId.id,
          displayName: 'test',
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          isActive: true,
          createdAt: entity.createdAt,
        },
      },
      {
        input: {
          id: entity.userId.id,
          displayName: 'test',
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          isActive: true,
          createdAt: entity.createdAt,
        },
      },
      {
        input: {
          id: entity.userId.id,
          displayName: 'test',
          isActive: false,
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          isActive: false,
          createdAt: entity.createdAt,
        },
      },
      {
        input: {
          id: entity.userId.id,
          displayName: 'test',
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          isActive: false,
          createdAt: entity.createdAt,
        },
      },
      {
        input: {
          id: entity.userId.id,
          displayName: 'test',
          isActive: true,
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          isActive: true,
          createdAt: entity.createdAt,
        },
      },
      {
        input: {
          id: entity.userId.id,
          displayName: 'test',
          isActive: false,
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          isActive: false,
          createdAt: entity.createdAt,
        },
      },
    ];

    for (const i of arrange) {
      output = await useCase.execute({
        id: i.input.id,
        ...(i.input.displayName && { displayName: i.input.displayName }),
        ...('isActive' in i.input && { isActive: i.input.isActive }),
      });
      const entityUpdated = await repository.findById(new Uuid(i.input.id));
      console.log(entityUpdated!.toJSON());
      expect(output).toStrictEqual({
        id: entity.userId.id,
        displayName: i.expected.displayName,
        dustBalance: entity.dustBalance,
        isActive: i.expected.isActive,
        createdAt: entityUpdated!.createdAt,
      });
      expect(entityUpdated!.toJSON()).toStrictEqual({
        userId: entity.userId.id,
        displayName: i.expected.displayName,
        dustBalance: entity.dustBalance,
        isActive: i.expected.isActive,
        createdAt: entityUpdated!.createdAt,
      });
    }
  });

  it('should NOT update a user dustBalance', async () => {
    const entity = User.fake().aUser().build();
    repository.insert(entity);

    let output = await useCase.execute({
      id: entity.userId.id,
      displayName: 'test',
    });
    expect(output).toStrictEqual({
      id: entity.userId.id,
      displayName: 'test',
      dustBalance: entity.dustBalance,
      isActive: true,
      createdAt: entity.createdAt,
    });

    type Arrange = {
      input: {
        id: string;
        displayName: string;
        dustBalance?: number;
        isActive?: boolean;
      };
      expected: {
        id: string;
        displayName: string;
        dustBalance: number;
        isActive: boolean;
        createdAt: Date;
      };
    };
    const arrange: Arrange[] = [
      {
        input: {
          id: entity.userId.id,
          displayName: 'test',
          dustBalance: 500,
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          dustBalance: entity.dustBalance,
          isActive: true,
          createdAt: entity.createdAt,
        },
      },
      {
        input: {
          id: entity.userId.id,
          displayName: 'test',
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          dustBalance: entity.dustBalance,
          isActive: true,
          createdAt: entity.createdAt,
        },
      },
      {
        input: {
          id: entity.userId.id,
          displayName: 'test',
          isActive: false,
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          dustBalance: entity.dustBalance,
          isActive: false,
          createdAt: entity.createdAt,
        },
      },
      {
        input: {
          id: entity.userId.id,
          displayName: 'test',
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          dustBalance: entity.dustBalance,
          isActive: false,
          createdAt: entity.createdAt,
        },
      },
      {
        input: {
          id: entity.userId.id,
          displayName: 'test',
          isActive: true,
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          dustBalance: entity.dustBalance,
          isActive: true,
          createdAt: entity.createdAt,
        },
      },
      {
        input: {
          id: entity.userId.id,
          displayName: 'test',
          dustBalance: 0,
          isActive: false,
        },
        expected: {
          id: entity.userId.id,
          displayName: 'test',
          dustBalance: entity.dustBalance,
          isActive: false,
          createdAt: entity.createdAt,
        },
      },
    ];

    for (const i of arrange) {
      output = await useCase.execute({
        id: i.input.id,
        ...(i.input.displayName && { displayName: i.input.displayName }),
        ...('dustBalance' in i.input && { dustBalance: i.input.dustBalance }),
        ...('isActive' in i.input && { isActive: i.input.isActive }),
      });
      const entityUpdated = await repository.findById(new Uuid(i.input.id));
      console.log(entityUpdated!.toJSON());
      expect(output).toStrictEqual({
        id: entity.userId.id,
        displayName: i.expected.displayName,
        dustBalance: entity.dustBalance,
        isActive: i.expected.isActive,
        createdAt: entityUpdated!.createdAt,
      });
      expect(entityUpdated!.toJSON()).toStrictEqual({
        userId: entity.userId.id,
        displayName: i.expected.displayName,
        dustBalance: entity.dustBalance,
        isActive: i.expected.isActive,
        createdAt: entityUpdated!.createdAt,
      });
    }
  });
});
