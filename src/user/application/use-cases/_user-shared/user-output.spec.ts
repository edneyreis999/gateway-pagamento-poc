import { User } from '../../../domain/user.entity';
import { UserOutputMapper } from './user-output';

describe('UserOutputMapper Unit Tests', () => {
  it('should convert a user in output', () => {
    const entity = User.create({
      displayName: 'Movie',
      dustBalance: 200,
      isActive: true,
    });
    const spyToJSON = jest.spyOn(entity, 'toJSON');
    const output = UserOutputMapper.toOutput(entity);
    expect(spyToJSON).toHaveBeenCalled();
    expect(output).toStrictEqual({
      id: entity.userId.id,
      displayName: 'Movie',
      dustBalance: 200,
      isActive: true,
      createdAt: entity.createdAt,
    });
  });
});
