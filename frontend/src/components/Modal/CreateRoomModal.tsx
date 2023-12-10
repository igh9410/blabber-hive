import { Button } from '@/shad-cn-ui/components/ui/button';
import { Input } from '@/shad-cn-ui/components/ui/input';

type CreateRoomModalProps = {
  isOpen: boolean;
  onClose: () => void;
  // Add any other props you might need
};

export function CreateRoomModal(props: Readonly<CreateRoomModalProps>) {
  if (!props.isOpen) return null; // If modal is closed, don't render anything

  return (
    <div
      aria-labelledby="modal-title"
      aria-modal="true"
      className="fixed z-10 inset-0 overflow-y-auto"
      role="dialog"
    >
      <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div
          aria-hidden="true"
          className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
        />
        <span
          aria-hidden="true"
          className="hidden sm:inline-block sm:align-middle sm:h-screen"
        ></span>
        <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div className="sm:flex sm:items-start">
              <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                <h3
                  className="text-lg leading-6 font-medium text-gray-900"
                  id="modal-title"
                >
                  Create a new room
                </h3>
                <div className="mt-2">
                  <Input
                    className="w-full bg-white shadow-none appearance-none pl-2 md:w-2/3 lg:w-2/3 dark:bg-gray-950"
                    placeholder="Room name..."
                    type="text"
                  />
                </div>
              </div>
            </div>
          </div>
          <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <Button className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm">
              Create
            </Button>
            <Button
              onClick={props.onClose}
              className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
            >
              Cancel
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
