import { IsBoolean, IsOptional, IsString, IsUUID, validateSync } from 'class-validator';

export type UpdateUserInputConstructorProps = {
  id: string;
  displayName?: string;
  isActive?: boolean;
};

export class UpdateUserInput {
  @IsUUID()
  id!: string;

  @IsString()
  @IsOptional()
  displayName?: string;

  @IsBoolean()
  @IsOptional()
  isActive?: boolean;

  constructor(props?: UpdateUserInputConstructorProps) {
    if (!props) return;
    this.id = props.id;
    props.displayName && (this.displayName = props.displayName);
    props.isActive !== null && props.isActive !== undefined && (this.isActive = props.isActive);
  }
}

export class ValidateUpdateUserInput {
  static validate(input: UpdateUserInputConstructorProps) {
    const addDustUserInput = new UpdateUserInput(input);
    const errors = validateSync(addDustUserInput);
    return errors;
  }
}
