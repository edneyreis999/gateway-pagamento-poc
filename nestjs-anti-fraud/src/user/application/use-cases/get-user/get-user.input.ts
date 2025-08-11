import { IsUUID, validateSync } from 'class-validator';

export type GetUserInputConstructorProps = {
  id: string;
};

export class GetUserInput {
  @IsUUID()
  id!: string;

  constructor(props?: GetUserInputConstructorProps) {
    if (!props) return;

    this.id = props.id;
  }
}

export class ValidateGetUserInput {
  static validate(input: GetUserInputConstructorProps) {
    const addDustUserInput = new GetUserInput(input);
    const errors = validateSync(addDustUserInput);
    return errors;
  }
}
