import { Notification } from './shared/domain/validators/notification';
declare global {
  namespace jest {
    interface Matchers<R> {
      notificationContainsErrorMessages: (
        received: Array<string | { [key: string]: string[] }>,
        expected?: Notification,
      ) => R;
    }
  }
}
