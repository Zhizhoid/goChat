//React imports
import * as React from 'react';

//Material UI imports
import Box from '@mui/material/Box';
import Drawer from '@mui/material/Drawer';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';

//Other imports

//Local imports
import RoomList from './Rooms/RoomList';
import ChatScreen from './Chat/ChatScreen';
import LoginDialog from './LoginDialog';

const drawerWidth = 240;
const backendIP = "http://localhost:8080"

// let testRoom = {
//     Name: "Test room 1",
//     Messages: [],
//     ID: 1,
//     //...
// }

// let testRoom2 = {
//     Name: "Test room 2",
//     Messages: [
//         {Text: "Foo", Author: "TheUser1"},
//         {Text: "Bar", Author: "TheUser2"},
//         {Text: "FooBar", Author: "TheUser1"},
//         {Text: "BarFoo", Author: "TheUser2"},
//     ],
//     ID: 2,
//     //...
// }

// let testRoom3 = {
//     Name: "Test room 3",
//     Messages: [
//         {Text: "Foo", Author: "TheUser1"},
//         {Text: "Bar", Author: "TheUser2"},
//         {Text: "FooBar", Author: "TheUser1"},
//         {Text: "BarFoo", Author: "TheUser2"},
//         {Text: "Foo", Author: "TheUser1"},
//         {Text: "Bar", Author: "TheUser2"},
//         {Text: "FooBar", Author: "TheUser1"},
//         {Text: "BarFoo", Author: "TheUser2"},
//         {Text: "Foo", Author: "TheUser1"},
//         {Text: "Bar", Author: "TheUser2"},
//         {Text: "FooBar", Author: "TheUser1"},
//         {Text: "BarFoo", Author: "TheUser2"},
//         {Text: "Foo", Author: "TheUser1"},
//         {Text: "Bar", Author: "TheUser2"},
//         {Text: "FooBar", Author: "TheUser1"},
//         {Text: "BarFoo", Author: "TheUser2"},
//         {Text: "Foo", Author: "TheUser1"},
//         {Text: "Bar", Author: "TheUser2"},
//         {Text: "FooBar", Author: "TheUser1"},
//         {Text: "BarFoo", Author: "TheUser2"},
//         {Text: "Foo", Author: "TheUser1"},
//         {Text: "Bar", Author: "TheUser2"},
//         {Text: "FooBar", Author: "TheUser1"},
//         {Text: "BarFoo", Author: "TheUser2"},
//         {Text: "Foo", Author: "TheUser1"},
//         {Text: "Bar", Author: "TheUser2"},
//         {Text: "FooBar", Author: "TheUser1"},
//         {Text: "BarFoo", Author: "TheUser2"},
//     ],
//     ID: 3,
//     //...
// }

const emptyRoom = {
    Name: "",
    Messages: [],
    ID: 0,
    //...
}

export default function MainScreen() {
    const [token, setToken] = React.useState("");

    const [roomList, setRoomList] = React.useState([]);
    const [activeRoom, setActiveRoom] = React.useState(emptyRoom);

    function initEmptyRooms(roomNames) {
        const rooms = [];
        for(let i = 0; i < roomNames.length; i++) {
            rooms.push({
                Name:roomNames[i],
                Messages: [],
                ID: i,
            });
        }
        return rooms;
    }

    function updateRoomList() {
        // setRoomList([testRoom, testRoom2, testRoom3]);
        let actn = {
            action: "read",
            object: "user",
            data: {
                token: token,
                getRoomList: true,
            },
        };

        fetch(backendIP.concat("/"), {
			method: 'POST', 
			mode: 'cors', 
			cache: 'no-cache', 
			credentials: 'same-origin', 
			headers: {
				'Content-Type': 'application/json'
			},
			redirect: 'follow',
			referrerPolicy: 'no-referrer', 
			body: JSON.stringify(actn),
		}).then(resp => {
			if (!resp.ok) {
				console.error("Error while reading rooms");
			}

			return resp.json();
		}).then(data => {
			if(!data.success) {
				console.error(data.status);
                return
			}

            setRoomList(initEmptyRooms(data.read.user.rooms));
		});
    }

    React.useEffect(() => {
        if(token === "") return;
        updateRoomList();
    });

    return (
        <Box sx={{ display: 'flex' }} height="100%">  {/*container for everything*/} 

            {/*AppBar is the blue bar with the title on top*/}
            <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
                <Toolbar>
                    <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
                        The Go Chat: {activeRoom.Name}
                    </Typography>
                    
                    <LoginDialog backendIP={backendIP} setToken={setToken}/>
                </Toolbar>
            </AppBar>

            {/*Drawer is that thing on the left side*/}
            <Drawer
                variant="permanent"
                sx={{
                    width: drawerWidth,
                    flexShrink: 0,
                    [`& .MuiDrawer-paper`]: { width: drawerWidth, boxSizing: 'border-box' },
                }}
            >
                <Toolbar />
                <RoomList activeRoom={activeRoom} setActiveRoom={setActiveRoom} roomList={roomList}/>
            </Drawer>

            {/*This is the window with the chat*/}
            <ChatScreen activeRoom={activeRoom} setActiveRoom={setActiveRoom}/>
        </Box>
    );
}
