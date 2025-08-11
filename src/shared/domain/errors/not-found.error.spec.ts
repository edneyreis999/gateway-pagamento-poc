import { NotFoundError } from './not-found.error';
import { Entity } from '../entity';
import { ValueObject } from '../value-object';

class TestEntity extends Entity {
  get entityId(): ValueObject {
    throw new Error('Method not implemented.');
  }
  toJSON() {
    throw new Error('Method not implemented.');
  }
}

describe('NotFoundError', () => {
  it('should format error message correctly with single ID', () => {
    const error = new NotFoundError('123', TestEntity);
    expect(error.message).toEqual('TestEntity Not Found using ID 123');
  });

  it('should format error message correctly with multiple IDs', () => {
    const error = new NotFoundError(['123', '456'], TestEntity);
    expect(error.message).toEqual('TestEntity Not Found using ID 123, 456');
  });
});
