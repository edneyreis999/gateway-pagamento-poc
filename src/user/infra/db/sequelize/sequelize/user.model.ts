import { Column, DataType, Model, PrimaryKey, Table } from 'sequelize-typescript';

export type UserModelProps = {
  userId: string;
  displayName: string;
  dustBalance: number | null;
  isActive: boolean;
  createdAt: Date;
};

@Table({ tableName: 'users', timestamps: false })
export class UserModel extends Model<UserModelProps> {
  @PrimaryKey
  @Column({ type: DataType.UUID })
  declare userId: string;

  @Column({ allowNull: false, type: DataType.STRING(255) })
  declare displayName: string;

  @Column({ allowNull: false, type: DataType.NUMBER })
  declare dustBalance: number | null;

  @Column({ allowNull: false, type: DataType.BOOLEAN })
  declare isActive: boolean;

  @Column({ allowNull: false, type: DataType.DATE(3) })
  declare createdAt: Date;
}
