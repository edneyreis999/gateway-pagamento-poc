import { IsNumber, IsUUID, Min, validateSync } from 'class-validator';

export type SubtractDustUserInputConstructorProps = {
  id: string;
  dust: number;
};

export class SubtractDustUserInput {
  @IsUUID()
  id!: string;

  @IsNumber()
  @Min(0)
  dust!: number;

  constructor(props?: SubtractDustUserInputConstructorProps) {
    if (!props) return;

    this.id = props.id;
    this.dust = props.dust;
  }
}

export class ValidateSubtractDustUserInput {
  static validate(input: SubtractDustUserInputConstructorProps) {
    const addDustUserInput = new SubtractDustUserInput(input);
    const errors = validateSync(addDustUserInput);
    return errors;
  }
}
