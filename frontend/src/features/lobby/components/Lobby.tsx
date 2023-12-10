import { CreateRoomModal } from '@components/Modal';
import { useChatRoomList } from '@hooks';
import { Button } from '@/shad-cn-ui/components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shad-cn-ui/components/ui/table';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

export function Lobby() {
  const { chatRooms, isLoading, error } = useChatRoomList();
  const navigate = useNavigate();

  const [isModalOpen, setIsModalOpen] = useState<boolean>(false); // Modal State for Create Room

  if (isLoading) return <div>Loading chat rooms...</div>;
  if (error) return <div>Error loading chat rooms: {error.message}</div>;

  function handleJoinRoom(id: string) {
    navigate(`/chats/${id}`);
  }

  const openModal = () => {
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
  };
  return (
    <>
      <main className="flex flex-1 flex-col gap-4 p-4 md:gap-8 md:p-6">
        <div className="flex items-center">
          <h1 className="font-semibold text-lg md:text-2xl">Chat Rooms</h1>
          <Button className="ml-auto" size="sm" onClick={openModal}>
            Create Room
          </Button>
        </div>
        <div className="border shadow-sm rounded-lg">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="max-w-[150px]">Name</TableHead>

                <TableHead>Join</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {chatRooms?.map((room) => (
                <TableRow key={room.id}>
                  <TableCell className="font-medium">
                    {room.name || 'Unnamed Room'}
                  </TableCell>
                  <TableCell>
                    <Button size="sm" onClick={() => handleJoinRoom(room.id)}>
                      Join
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </main>

      {isModalOpen && (
        <CreateRoomModal isOpen={isModalOpen} onClose={closeModal} />
      )}
    </>
  );
}
