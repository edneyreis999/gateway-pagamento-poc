import {
  IsBoolean,
  IsNotEmpty,
  IsNumber,
  IsOptional,
  IsString,
  Max,
  Min,
  validateSync,
} from 'class-validator';

export type CreateUserInputConstructorProps = {
  displayName: string;
  dustBalance?: number;
  isActive?: boolean;
};

export class CreateUserInput {
  @IsString()
  @IsNotEmpty()
  displayName!: string;

  @IsNumber()
  @IsOptional()
  @Min(0)
  @Max(999999)
  dustBalance?: number;

  @IsBoolean()
  @IsOptional()
  isActive?: boolean;

  constructor(props: CreateUserInputConstructorProps) {
    if (!props) return;

    this.displayName = props.displayName;
    props.dustBalance && (this.dustBalance = props.dustBalance);
    props.isActive && (this.isActive = props.isActive);
  }
}

export class ValidateCreateUserInput {
  static validate(input: CreateUserInputConstructorProps) {
    const createUserInput = new CreateUserInput(input);
    const errors = validateSync(createUserInput);
    return errors;
  }
}
