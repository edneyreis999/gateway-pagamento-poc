import { IsUUID, validateSync } from 'class-validator';

export type DeleteUserInputConstructorProps = {
  id: string;
};

export class DeleteUserInput {
  @IsUUID()
  id!: string;

  constructor(props: DeleteUserInputConstructorProps) {
    if (!props) return;
    this.id = props.id;
  }
}

export class ValidateDeleteUserInput {
  static validate(input: DeleteUserInputConstructorProps) {
    const deleteUserInput = new DeleteUserInput(input);
    const errors = validateSync(deleteUserInput);
    return errors;
  }
}
