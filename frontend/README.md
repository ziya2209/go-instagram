# Instagram Clone Frontend

A modern, responsive frontend for the Instagram Clone backend built with vanilla HTML, CSS, and JavaScript.

## Features

### 🔐 Authentication
- **Login**: Secure user authentication with JWT tokens
- **Signup**: User registration with validation
- **Auto-login**: Persistent sessions using localStorage

### 📱 Core Functionality
- **User Management**: View all users, follow/unfollow functionality
- **Post Creation**: Create posts with images and captions
- **Social Features**: Like posts, add comments (including nested comments)
- **Home Feed**: Centralized view for user content
- **Followers**: View people you follow

### 🎨 UI/UX Features
- **Responsive Design**: Mobile-first approach, works on all devices
- **Modern UI**: Instagram-inspired design with gradients and animations
- **Toast Notifications**: Real-time feedback for user actions
- **Loading States**: Visual feedback during API calls
- **Modal Windows**: Elegant popups for comments and follow actions
- **Smooth Animations**: Fade-in effects and hover animations

## File Structure

```
frontend/
├── index.html      # Main HTML structure
├── styles.css      # All styling and responsive design
├── script.js       # JavaScript functionality and API integration
└── README.md       # This file
```

## Setup Instructions

### 1. Backend Setup
Make sure your Go backend is running on `http://localhost:8080`. You can start it using:

```bash
# From the root directory of your Go project
go run main.go

# Or using Docker
make up

# Or using Make
make redeploy
```

### 2. Frontend Setup
1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Serve the files using any HTTP server:

   **Option A: Using Python (if installed)**
   ```bash
   # Python 3
   python -m http.server 3000
   
   # Python 2
   python -S SimpleHTTPServer 3000
   ```

   **Option B: Using Node.js (if installed)**
   ```bash
   npx http-server -p 3000
   ```

   **Option C: Using PHP (if installed)**
   ```bash
   php -S localhost:3000
   ```

   **Option D: Using VS Code Live Server**
   - Install the "Live Server" extension
   - Right-click on `index.html`
   - Select "Open with Live Server"

3. Open your browser and navigate to `http://localhost:3000`

## API Endpoints Used

The frontend integrates with the following backend endpoints:

### Public Endpoints
- `GET /health` - Health check
- `POST /createAccount` - User registration
- `POST /login` - User authentication

### Protected Endpoints (require JWT token)
- `GET /users` - Get all users
- `POST /follow` - Follow a user
- `GET /followers` - Get user's followers
- `POST /createPost` - Create a new post
- `POST /likePost` - Like a post
- `POST /addComment/post` - Add comment to post
- `POST /addComment/comment` - Add comment to comment (nested)
- `GET /post/comments` - Get post comments
- `GET /home` - Home feed (placeholder)

## Usage Guide

### Getting Started
1. **First Time Users**: Click "Sign up" to create an account
2. **Existing Users**: Use your username and password to log in

### Creating Content
1. Click the **+** button in the navigation
2. Enter an image URL and caption
3. Click "Share Post"

### Social Features
1. **Following Users**: 
   - Go to "Users" tab
   - Click "Follow" on any user card
   - Or use the heart icon and search for specific users

2. **Viewing Followers**:
   - Click the heart icon in navigation
   - See all users you're following

3. **Commenting** (when posts are available):
   - Click comment icon on any post
   - Add your comment in the modal

4. **Liking Posts**:
   - Click the heart icon on any post

### Navigation
- **Home** (🏠): Main feed (placeholder for now)
- **Create** (➕): Create new posts
- **Users** (👥): Browse and follow users
- **Following** (❤️): View your followers
- **Logout** (🚪): Sign out of the application

## Browser Support
- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+

## CORS Configuration
The frontend is configured to work with the backend's CORS settings. The backend allows:
- All origins (`*`)
- Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
- Headers: Origin, Content-Type, Accept, Authorization

## Responsive Breakpoints
- **Desktop**: 1200px+
- **Tablet**: 768px - 1199px
- **Mobile**: Below 768px

## Key Features Explained

### Authentication Flow
1. JWT tokens are stored in localStorage
2. Automatic token inclusion in API requests
3. Persistent login across browser sessions
4. Secure logout clears all stored data

### Error Handling
- Network errors show user-friendly messages
- API errors are parsed and displayed
- Loading states prevent multiple submissions
- Form validation prevents invalid data

### Performance Features
- Efficient DOM manipulation
- Lazy loading of user data
- Minimal API calls
- CSS transitions for smooth UX

## Customization

### Styling
- Modify `styles.css` to change colors, fonts, or layout
- CSS custom properties make theming easy
- Responsive grid system adapts to content

### Functionality
- Add new API endpoints in `script.js`
- Extend UI components in `index.html`
- Modify API base URL in the configuration section

## Troubleshooting

### Common Issues

1. **"Unable to connect to server"**
   - Ensure backend is running on port 8080
   - Check if Docker containers are up
   - Verify CORS configuration

2. **Login not working**
   - Check if user exists in database
   - Verify password is correct
   - Check browser console for errors

3. **Features not loading**
   - Ensure you're logged in
   - Check JWT token in localStorage
   - Verify API endpoints are accessible

### Development Tips
- Use browser developer tools to debug
- Check Network tab for API call details
- Monitor Console for JavaScript errors
- Use Application tab to view localStorage

## Future Enhancements
- Real-time notifications
- Image upload functionality
- User profiles with edit capabilities
- Advanced search and filtering
- Dark mode theme
- Progressive Web App (PWA) features

## Contributing
Feel free to submit issues and enhancement requests!
