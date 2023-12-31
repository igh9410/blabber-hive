import { WEBSOCKET_URL, authTokenKey } from '@config';
import { MessageType, ServerMessageType } from '@features/chat';
import { useAuth } from '@providers';
import { useChatMessageStore } from '@stores';
import { useQueryClient } from '@tanstack/react-query';
import { useEffect, useRef, useState } from 'react';

export const useWebSocketConnection = (chatRoomId: string) => {
  const queryClient = useQueryClient();
  const { session } = useAuth();
  const addMessage = useChatMessageStore((state) => state.addMessage);

  const [retryCount, setRetryCount] = useState(0);
  const [isConnected, setIsConnected] = useState(false);

  const webSocket = useRef<WebSocket | null>(null);

  const sendMessage = (message: string) => {
    if (webSocket.current) {
      webSocket.current.send(message); // serialize the object to a JSON string
    }
  };

  useEffect(() => {
    const supabaseData = localStorage.getItem(authTokenKey);

    let accessToken = '';

    if (supabaseData) {
      const parsedData = JSON.parse(supabaseData);
      accessToken = parsedData.access_token;
    } else {
      console.error('No Supabase data found in Local Storage');
    }

    const webSocketURL =
      WEBSOCKET_URL + '/chats/' + chatRoomId + '?token=' + accessToken;

    webSocket.current = new WebSocket(webSocketURL);

    const handleWebSocketOpen = () => {
      console.log('WebSocket connection opened');
      setIsConnected(true);
    };

    const handleWebSocketOnMessage = (event: MessageEvent) => {
      const serverMessage: ServerMessageType = JSON.parse(event.data);
      console.log('Server message = ', serverMessage);
      // Convert the snake_case properties from the server message to camelCase for the MessageType
      // Explicitly define the type of messageSenderId as 'sent' | 'received'
      let messageSenderId: 'sent' | 'received';

      if (session?.user?.id === serverMessage.sender_id) {
        messageSenderId = 'sent';
      } else {
        messageSenderId = 'received';
      }

      const receivedMessage: MessageType = {
        sender: messageSenderId,
        senderID: serverMessage.sender_id,
        content: serverMessage.content,
        createdAt: new Date(serverMessage.created_at),
        img: serverMessage.media_url || undefined,
      };

      addMessage(receivedMessage);
    };

    const reconnetWebSocket = (error: Event) => {
      if (retryCount < 3 && !isConnected) {
        console.log(
          'WebSocket connection closed, trying reconnect..:',
          error.type
        );
        setRetryCount((prevRetryCount) => prevRetryCount + 1);
        setTimeout(() => {
          if (webSocket.current) {
            webSocket.current.close();
          }
          webSocket.current = new WebSocket(webSocketURL);
        }, 2000);
        console.log('Trying to reconnect, Retry Count = ', retryCount);
      }
    };

    const handleWebSocketOnError = (error: Event) => {
      console.error(`Socket encountered error: Closing socket`, error.type);
      reconnetWebSocket(error);
      if (webSocket.current) {
        webSocket.current.close();
      }
    };

    webSocket.current!.addEventListener('open', handleWebSocketOpen); // Non-null assertion here
    webSocket.current!.addEventListener('message', handleWebSocketOnMessage); // Non-null assertion here
    webSocket.current!.addEventListener('close', reconnetWebSocket); // Non-null assertion here
    webSocket.current!.addEventListener('error', handleWebSocketOnError); // Non-null assertion here

    return () => {
      console.log('Performing Cleanup function');
      if (webSocket.current) {
        webSocket.current.removeEventListener('open', handleWebSocketOpen);
        webSocket.current.removeEventListener(
          'message',
          handleWebSocketOnMessage
        );
        webSocket.current.removeEventListener('close', reconnetWebSocket);
        webSocket.current.removeEventListener('error', handleWebSocketOnError);
        webSocket.current.close();
      }
      setIsConnected(false);
      console.log('WebSocket connection closed');
    };
  }, [queryClient, retryCount]); // Added dependencies

  return { sendMessage };
};
