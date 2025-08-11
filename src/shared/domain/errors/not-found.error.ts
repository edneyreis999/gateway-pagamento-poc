import { Entity } from '../entity';
import type { ValueObject } from '../value-object';

export class NotFoundError extends Error {
  constructor(id: any[] | string | ValueObject, entityClass: new (...args: any[]) => Entity) {
    const idsMessage = Array.isArray(id) ? id.join(', ') : id;
    super(`${entityClass.name} Not Found using ID ${idsMessage}`);
    this.name = 'NotFoundError';
  }
}
