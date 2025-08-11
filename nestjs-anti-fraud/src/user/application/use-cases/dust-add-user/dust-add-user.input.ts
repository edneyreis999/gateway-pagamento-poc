import { IsNumber, IsUUID, Min, validateSync } from 'class-validator';

export type AddDustUserInputConstructorProps = {
  id: string;
  dust: number;
};

export class AddDustUserInput {
  @IsUUID()
  id!: string;

  @IsNumber()
  @Min(0)
  dust!: number;

  constructor(props?: AddDustUserInputConstructorProps) {
    if (!props) return;

    this.id = props.id;
    this.dust = props.dust;
  }
}

export class ValidateAddDustUserInput {
  static validate(input: AddDustUserInputConstructorProps) {
    const addDustUserInput = new AddDustUserInput(input);
    const errors = validateSync(addDustUserInput);
    return errors;
  }
}
