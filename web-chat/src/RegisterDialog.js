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


export default function LoginDialog(props) {
	const [open, setOpen] = React.useState(false);

    const [name, setName] = React.useState("");
	const [username, setUsername] = React.useState("");
	const [password, setPassword] = React.useState("");

    function nameChange(event) {
        setName(event.target.value);
    }

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

	function handleRegistration() {
		let actn = {
			action: "create",
			object: "user",
			data: {
				username: username,
				password: password,
                name: name,
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
			if (!resp.ok) {
				alert("Error occured during registration");
			}

			return resp.json();
		}).then(data => {
			console.log(data);
			
			if(!data.success) {
				alert("Registration failed\nMaybe the username or the password are already used or empty");
			} else {
				alert("Registered successfully")
			}

			setOpen(false);
		});
	}

	return (
		<>
			<Button variant="standard" onClick={handleOpen}>
				Register
			</Button>
			<Dialog open={open} onClose={handleClose}>
				<DialogTitle>Registration</DialogTitle>
				<DialogContent>
					<DialogContentText>
						Enter your credentials
					</DialogContentText>
                    <TextField
						autoFocus
						margin="dense"
						label="Name"
						type="text"
						fullWidth
						variant="standard"
						value={name}
						onChange={nameChange}
					/>
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
					<Button onClick={handleRegistration}>Register</Button>
				</DialogActions>
			</Dialog>
		</>
	);
}

LoginDialog.propTypes = {
    backendIP: PropTypes.any.isRequired,
};