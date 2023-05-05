//React imports
import * as React from 'react';

//Material UI imports
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';

//Other imports
import PropTypes from 'prop-types';

//Local imports
import RegisterDialog from './RegisterDialog';

export default function LoginDialog(props) {
	const [open, setOpen] = React.useState(true);
	const [username, setUsername] = React.useState("");
	const [password, setPassword] = React.useState("");
	const [loggedIn, setLoggedIn] = React.useState(false);

	function usernameChange(event) {
		setUsername(event.target.value);
	};

	function passwordChange(event) {
		setPassword(event.target.value);
	};

	function handleOpen() {
		setOpen(true);
	};

	function handleClose() {
		setOpen(false);
	};

	function handleLogin() {
		let actn = {
			action: "login",
			object: "user",
			data: {
				username: username,
				password: password,
			}
		}

		fetch(props.backendIP.concat("/"), {
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
			//The place where you should check if request was successfull and read info about response like headers
			if (!resp.ok) {
				alert("Error occured during login");
			}

			return resp.json();
		}).then(data => {
			// console.log(data);
			
			if(!data.success) {
				alert("Login failed\nMaybe the username or the password is invalid");
				setOpen(false);
				return;
			}

			setLoggedIn(true);
			props.setToken(data.token);

			setOpen(false);
		});
	}

	return (
		<>
			{loggedIn ? null : <RegisterDialog backendIP={props.backendIP}/>}

			<Button variant="standard" onClick={handleOpen}>
				{loggedIn ? username : "Log in"}
			</Button>
			<Dialog open={open} onClose={handleClose}>
				<DialogTitle>Login</DialogTitle>
				<DialogContent>
					<DialogContentText>
						Enter your credentials
					</DialogContentText>
					<TextField
						autoFocus
						margin="dense"
						label="Username"
						type="text"
						fullWidth
						variant="standard"
						value={username}
						onChange={usernameChange}
					/>
					<TextField
						margin="dense"
						label="Password"
						type="password"
						fullWidth
						variant="standard"
						value={password}
						onChange={passwordChange}
					/>
				</DialogContent>
				<DialogActions>
					<Button onClick={handleClose}>Cancel</Button>
					<Button onClick={handleLogin}>Log in</Button>
				</DialogActions>
			</Dialog>
		</>
	);
}

LoginDialog.propTypes = {
    backendIP: PropTypes.any.isRequired,
	setToken: PropTypes.func.isRequired
};