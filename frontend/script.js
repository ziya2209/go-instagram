// API Configuration
const API_BASE_URL = 'http://localhost:8080';

// Cookie utility functions
function setCookie(name, value, days = 7) {
    const expires = new Date();
    expires.setTime(expires.getTime() + (days * 24 * 60 * 60 * 1000));
    document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/;SameSite=Lax`;
}

function getCookie(name) {
    const nameEQ = name + "=";
    const ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == ' ') c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
}

function deleteCookie(name) {
    document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:01 GMT;path=/;SameSite=Lax`;
}

// Get auth data from cookies (fallback to localStorage for backward compatibility)
let authToken = getCookie('authToken') || localStorage.getItem('authToken');
let currentUser = getCookie('currentUser') || localStorage.getItem('currentUser');

// DOM Elements
const loginSection = document.getElementById('loginSection');
const signupSection = document.getElementById('signupSection');
const mainContent = document.getElementById('mainContent');
const navbar = document.getElementById('navbar');
const loading = document.getElementById('loading');
const toast = document.getElementById('toast');

// Navigation elements
const homeBtn = document.getElementById('homeBtn');
const createPostBtn = document.getElementById('createPostBtn');
const usersBtn = document.getElementById('usersBtn');
const followersBtn = document.getElementById('followersBtn');
const logoutBtn = document.getElementById('logoutBtn');

// Content sections
const homeSection = document.getElementById('homeSection');
const createPostSection = document.getElementById('createPostSection');
const usersSection = document.getElementById('usersSection');
const followersSection = document.getElementById('followersSection');

// Modals
const postModal = document.getElementById('postModal');
const commentModal = document.getElementById('commentModal');
const followModal = document.getElementById('followModal');

// Utility Functions
function showLoading() {
    loading.style.display = 'block';
}

function hideLoading() {
    loading.style.display = 'none';
}

function showToast(message, type = 'success') {
    toast.textContent = message;
    toast.className = `toast ${type} show`;
    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
}

// API Functions
async function apiCall(endpoint, options = {}) {
    showLoading();
    try {
        const config = {
            ...options
        };

        // Set headers only when needed
        if (!config.headers) {
            config.headers = {};
        }

        // Add Content-Type for POST requests with body
        if (config.method === 'POST' && config.body) {
            config.headers['Content-Type'] = 'application/json';
        }

        // Add Authorization header for protected endpoints (not for login/signup)
        // Update: Always add Authorization header if token is available, except for login/signup
        // Refresh token from cookies in case it was updated
        const currentToken = getCookie('authToken') || localStorage.getItem('authToken');
        if (currentToken && !endpoint.includes('/login') && !endpoint.includes('/createAccount')) {
            authToken = currentToken; // Update the global variable
            config.headers['Authorization'] = `Bearer ${currentToken}`;
        }

        console.log('Making API call:', endpoint, 'with headers:', config.headers); // Debug log

        const response = await fetch(`${API_BASE_URL}${endpoint}`, config);
        
        // Handle non-JSON responses
        let data;
        const contentType = response.headers.get('content-type');
        if (contentType && contentType.includes('application/json')) {
            data = await response.json();
        } else {
            const textResponse = await response.text();
            try {
                // Try to parse as JSON even if content-type is wrong
                data = JSON.parse(textResponse);
            } catch (e) {
                data = { message: textResponse };
            }
        }

        if (!response.ok) {
            throw new Error(data.error || data.message || `HTTP ${response.status}: ${response.statusText}`);
        }

        return data;
    } catch (error) {
        showToast(error.message, 'error');
        throw error;
    } finally {
        hideLoading();
    }
}

// Authentication Functions
async function login(username, password) {
    try {
        const data = await apiCall('/login', {
            method: 'POST',
            body: JSON.stringify({ username, password })
        });

        // Store token in cookies (preferred) and localStorage (fallback)
        authToken = data.token;
        currentUser = username;
        
        setCookie('authToken', authToken, 7); // Store for 7 days
        setCookie('currentUser', currentUser, 7);
        
        // Also store in localStorage for backward compatibility
        localStorage.setItem('authToken', authToken);
        localStorage.setItem('currentUser', currentUser);

        console.log('Login successful, token stored:', authToken.substring(0, 20) + '...'); // Debug log

        showMainApp();
        showToast('Logged in successfully!');
    } catch (error) {
        console.error('Login failed:', error);
    }
}

async function signup(userData) {
    try {
        await apiCall('/createAccount', {
            method: 'POST',
            body: JSON.stringify(userData)
        });

        showToast('Account created successfully! Please log in.');
        showLogin();
    } catch (error) {
        console.error('Signup failed:', error);
    }
}

function logout() {
    authToken = null;
    currentUser = null;
    
    // Delete cookies
    deleteCookie('authToken');
    deleteCookie('currentUser');
    
    // Also clear localStorage for backward compatibility
    localStorage.removeItem('authToken');
    localStorage.removeItem('currentUser');
    
    hideMainApp();
    showLogin();
    showToast('Logged out successfully!');
}

// App Navigation
function showMainApp() {
    loginSection.style.display = 'none';
    signupSection.style.display = 'none';
    mainContent.style.display = 'block';
    navbar.style.display = 'block';
    
    showHome();
}

function hideMainApp() {
    mainContent.style.display = 'none';
    navbar.style.display = 'none';
}

function showLogin() {
    signupSection.style.display = 'none';
    loginSection.style.display = 'flex';
}

function showSignup() {
    loginSection.style.display = 'none';
    signupSection.style.display = 'flex';
}

function showSection(sectionToShow) {
    [homeSection, createPostSection, usersSection, followersSection].forEach(section => {
        section.style.display = 'none';
    });
    sectionToShow.style.display = 'block';
}

function showHome() {
    showSection(homeSection);
    loadHomeFeed();
}

function showCreatePost() {
    showSection(createPostSection);
}

function showUsers() {
    showSection(usersSection);
    loadAllUsers();
}

function showFollowers() {
    showSection(followersSection);
    loadFollowers();
}

// Data Loading Functions
async function loadHomeFeed() {
    try {
        // Since ShowHomePage is not implemented, we'll create a mock home feed
        const postsContainer = document.getElementById('posts');
        postsContainer.innerHTML = '<div class="text-center"><p>Home feed coming soon! Create some posts and follow users to see content here.</p></div>';
        
        // For now, we'll show a message that the feature is coming soon
        showToast('Home feed feature is being developed. Check back soon!', 'info');
    } catch (error) {
        console.error('Failed to load home feed:', error);
    }
}

async function loadAllUsers() {
    try {
        const users = await apiCall('/users');
        const usersList = document.getElementById('usersList');
        
        usersList.innerHTML = users.map(user => `
            <div class="user-card fade-in">
                <div class="user-avatar">${user.username.charAt(0).toUpperCase()}</div>
                <div class="user-info">
                    <h3>${user.username}</h3>
                    <p>${user.email}</p>
                    <div class="user-bio">${user.bio || 'No bio available'}</div>
                    <p>Age: ${user.age}</p>
                    <button class="btn btn-primary btn-small" onclick="openFollowModal('${user.username}')">
                        Follow
                    </button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Failed to load users:', error);
    }
}

async function loadFollowers() {
    try {
        const data = await apiCall('/followers');
        const followersList = document.getElementById('followersList');
        
        if (data.followers && data.followers.length > 0) {
            followersList.innerHTML = data.followers.map(username => `
                <div class="follower-item fade-in">
                    <div class="follower-avatar">${username.charAt(0).toUpperCase()}</div>
                    <div class="follower-name">${username}</div>
                </div>
            `).join('');
        } else {
            followersList.innerHTML = '<div class="text-center"><p>You are not following anyone yet. Go to the Users section to follow people!</p></div>';
        }
    } catch (error) {
        console.error('Failed to load followers:', error);
    }
}

// Post Functions
async function createPost(url, caption) {
    try {
        await apiCall('/createPost', {
            method: 'POST',
            body: JSON.stringify({ url, caption })
        });

        showToast('Post created successfully!');
        document.getElementById('createPostForm').reset();
        showHome();
    } catch (error) {
        console.error('Failed to create post:', error);
    }
}

async function likePost(postId) {
    try {
        await apiCall('/likePost', {
            method: 'POST',
            body: JSON.stringify({ post_id: postId })
        });

        showToast('Post liked!');
    } catch (error) {
        console.error('Failed to like post:', error);
    }
}

// Comment Functions
async function loadPostComments(postId) {
    try {
        const comments = await apiCall(`/post/comments?postId=${postId}`);
        return comments;
    } catch (error) {
        console.error('Failed to load comments:', error);
        return [];
    }
}

async function addComment(parentId, parentType, comment) {
    try {
        await apiCall('/addComment/post', {
            method: 'POST',
            body: JSON.stringify({
                parent_id: parentId,
                parent_type: parentType,
                comment: comment
            })
        });

        showToast('Comment added successfully!');
        closeModal(commentModal);
    } catch (error) {
        console.error('Failed to add comment:', error);
    }
}

// Follow Functions
async function followUser(username) {
    try {
        await apiCall('/follow', {
            method: 'POST',
            body: JSON.stringify({ username })
        });

        showToast(`You are now following ${username}!`);
        closeModal(followModal);
        if (followersSection.style.display !== 'none') {
            loadFollowers();
        }
    } catch (error) {
        console.error('Failed to follow user:', error);
    }
}

// Modal Functions
function openModal(modal) {
    modal.style.display = 'block';
}

function closeModal(modal) {
    modal.style.display = 'none';
}

function openFollowModal(username = '') {
    document.getElementById('followUsername').value = username;
    openModal(followModal);
}

function openCommentModal(parentId, parentType) {
    document.getElementById('commentParentId').value = parentId;
    document.getElementById('commentParentType').value = parentType;
    document.getElementById('commentText').value = '';
    openModal(commentModal);
}

// Event Listeners
document.addEventListener('DOMContentLoaded', function() {
    // Refresh token from cookies on page load
    authToken = getCookie('authToken') || localStorage.getItem('authToken');
    currentUser = getCookie('currentUser') || localStorage.getItem('currentUser');
    
    console.log('App initialized with token:', authToken ? authToken.substring(0, 20) + '...' : 'No token'); // Debug log
    
    // Check if user is already logged in
    if (authToken && currentUser) {
        showMainApp();
    } else {
        showLogin();
    }

    // Auth form handlers
    document.getElementById('showSignup').addEventListener('click', (e) => {
        e.preventDefault();
        showSignup();
    });

    document.getElementById('showLogin').addEventListener('click', (e) => {
        e.preventDefault();
        showLogin();
    });

    // Login form
    document.getElementById('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('loginUsername').value;
        const password = document.getElementById('loginPassword').value;
        await login(username, password);
    });

    // Signup form
    document.getElementById('signupForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const userData = {
            username: document.getElementById('signupUsername').value,
            email: document.getElementById('signupEmail').value,
            password: document.getElementById('signupPassword').value,
            age: parseInt(document.getElementById('signupAge').value),
            bio: document.getElementById('signupBio').value
        };
        await signup(userData);
    });

    // Navigation handlers
    homeBtn.addEventListener('click', showHome);
    createPostBtn.addEventListener('click', showCreatePost);
    usersBtn.addEventListener('click', showUsers);
    followersBtn.addEventListener('click', showFollowers);
    logoutBtn.addEventListener('click', logout);

    // Create post form
    document.getElementById('createPostForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const url = document.getElementById('postUrl').value;
        const caption = document.getElementById('postCaption').value;
        await createPost(url, caption);
    });

    // Comment form
    document.getElementById('commentForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const parentId = parseInt(document.getElementById('commentParentId').value);
        const parentType = document.getElementById('commentParentType').value;
        const comment = document.getElementById('commentText').value;
        await addComment(parentId, parentType, comment);
    });

    // Follow form
    document.getElementById('followForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('followUsername').value;
        await followUser(username);
    });

    // Modal close handlers
    document.querySelectorAll('.close').forEach(closeBtn => {
        closeBtn.addEventListener('click', (e) => {
            closeModal(e.target.closest('.modal'));
        });
    });

    // Close modal when clicking outside
    window.addEventListener('click', (e) => {
        if (e.target.classList.contains('modal')) {
            closeModal(e.target);
        }
    });

    // Keyboard shortcuts
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            document.querySelectorAll('.modal').forEach(modal => {
                if (modal.style.display === 'block') {
                    closeModal(modal);
                }
            });
        }
    });
});

// Health check function
async function checkHealth() {
    try {
        const response = await fetch(`${API_BASE_URL}/health`);
        const data = await response.json();
        console.log('API Health:', data);
        return data.status === 'all ok';
    } catch (error) {
        console.error('Health check failed:', error);
        return false;
    }
}

// Initialize app
checkHealth().then(isHealthy => {
    if (!isHealthy) {
        showToast('Unable to connect to server. Please make sure the backend is running on port 8080.', 'error');
    }
});

// Global functions for HTML onclick handlers
window.openFollowModal = openFollowModal;
window.openCommentModal = openCommentModal;
window.likePost = likePost;

// Debug function to check token status
window.debugTokens = function() {
    console.log('=== TOKEN DEBUG ===');
    console.log('Cookie authToken:', getCookie('authToken'));
    console.log('localStorage authToken:', localStorage.getItem('authToken'));
    console.log('Global authToken:', authToken);
    console.log('Cookie currentUser:', getCookie('currentUser'));
    console.log('localStorage currentUser:', localStorage.getItem('currentUser'));
    console.log('Global currentUser:', currentUser);
    console.log('===================');
};

// Test API call function
window.testUsersAPI = async function() {
    try {
        console.log('Testing /users API call...');
        const users = await apiCall('/users');
        console.log('Success:', users);
        return users;
    } catch (error) {
        console.error('Failed:', error);
        return null;
    }
};
