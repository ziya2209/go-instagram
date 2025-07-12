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
const myPostsBtn = document.getElementById('myPostsBtn');
const usersBtn = document.getElementById('usersBtn');
const followersBtn = document.getElementById('followersBtn');
const logoutBtn = document.getElementById('logoutBtn');

// Content sections
const homeSection = document.getElementById('homeSection');
const createPostSection = document.getElementById('createPostSection');
const myPostsSection = document.getElementById('myPostsSection');
const usersSection = document.getElementById('usersSection');
const followersSection = document.getElementById('followersSection');

// Modals
const postModal = document.getElementById('postModal');
const commentModal = document.getElementById('commentModal');
const followModal = document.getElementById('followModal');
const unfollowModal = document.getElementById('unfollowModal');

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
    if (!dateString || dateString === '') {
        return 'Just now';
    }

    try {
        const date = new Date(dateString);
        if (isNaN(date.getTime())) {
            return 'Just now';
        }

        const now = new Date();
        const diffMs = now - date;
        const diffMins = Math.floor(diffMs / 60000);
        const diffHours = Math.floor(diffMs / 3600000);
        const diffDays = Math.floor(diffMs / 86400000);

        if (diffMins < 1) return 'Just now';
        if (diffMins < 60) return `${diffMins}m ago`;
        if (diffHours < 24) return `${diffHours}h ago`;
        if (diffDays < 7) return `${diffDays}d ago`;

        return date.toLocaleDateString();
    } catch (error) {
        console.error('Error formatting date:', error);
        return 'Just now';
    }
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
    [homeSection, createPostSection, myPostsSection, usersSection, followersSection].forEach(section => {
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

function showMyPosts() {
    showSection(myPostsSection);
    loadMyPosts();
}

function showUsers() {
    showSection(usersSection);
    loadAllUsers();
}

function showFollowers() {
    showSection(followersSection);
    loadFollowers();
}

// Global variables for infinite scroll
let currentHomePage = 1;
let isLoadingHomeFeed = false;
let hasMoreHomePosts = true;
const homeFeedLimit = 10;

// Data Loading Functions
async function loadHomeFeed(page = 1, append = false) {
    if (isLoadingHomeFeed) return;

    try {
        isLoadingHomeFeed = true;
        console.log('Loading home feed...', { page, append }); // Debug log

        const data = await apiCall(`/home?page=${page}&limit=${homeFeedLimit}`);
        console.log('Home feed data received:', data); // Debug log

        const postsContainer = document.getElementById('posts');

        if (!append) {
            postsContainer.innerHTML = '';
            currentHomePage = 1;
        }

        // Use exact field names from HomepageResponse DTO
        if (data.posts && data.posts.length > 0) {
            const postsHtml = data.posts.map(post => createPostHTML(post)).join('');

            if (append) {
                postsContainer.insertAdjacentHTML('beforeend', postsHtml);
            } else {
                postsContainer.innerHTML = postsHtml;
            }

            // Update pagination state using exact field names
            hasMoreHomePosts = data.has_next || false;

            // Setup infinite scroll if not already done
            if (!append) {
                setupInfiniteScroll();
            }
        } else {
            if (!append) {
                postsContainer.innerHTML = '<div class="text-center"><p>No posts to show. Follow some users to see their posts in your feed!</p></div>';
            }
            hasMoreHomePosts = false;
        }
    } catch (error) {
        console.error('Failed to load home feed:', error);
        const postsContainer = document.getElementById('posts');
        if (!append) {
            postsContainer.innerHTML = '<div class="text-center"><p>Failed to load posts. Please try again.</p></div>';
        }
        showToast('Failed to load posts. Please try again.', 'error');
    } finally {
        isLoadingHomeFeed = false;
    }
}

// Helper function to create post HTML
function createPostHTML(post) {
    // Use exact field names from PostWithAuthor struct in dto package
    const authorUsername = post.author_username || 'Unknown';
    const postId = post.id;
    const postUrl = post.url || '';
    const postCaption = post.caption || '';
    const createdAt = post.created_at || '';
    const likeCount = post.like_count || 0;
    const commentCount = post.comment_count || 0;

    return `
        <div class="post-card fade-in" id="post-${postId}" data-post-id="${postId}">
            <div class="post-header">
                <div class="post-author">
                    <div class="author-avatar">${authorUsername.charAt(0).toUpperCase()}</div>
                    <div class="author-info">
                        <div class="author-name">${authorUsername}</div>
                        <div class="post-date">${formatDate(createdAt)}</div>
                    </div>
                </div>
            </div>
            <div class="post-image">
                <img src="${postUrl}" alt="Post image" onerror="this.src='https://via.placeholder.com/400x400?text=Image+Not+Found'">
            </div>
            <div class="post-content">
                <div class="post-stats">
                    <button class="action-btn like-btn" onclick="likePost(${postId})">
                        <i class="fas fa-heart"></i> <span class="like-count">${likeCount}</span>
                    </button>
                    <button class="action-btn comment-btn" onclick="toggleComments(${postId})" title="View and add comments">
                        <i class="fas fa-comment"></i> <span class="comment-count">${commentCount}</span> comments
                    </button>
                </div>
                <div class="post-caption"><strong>${authorUsername}</strong> ${postCaption}</div>

                <!-- Comments section (initially hidden) -->
                <div class="comments-section" id="comments-${postId}" style="display: none;">
                    <div class="comments-header">
                        <h4>Comments</h4>
                        <button class="close-comments-btn" onclick="toggleComments(${postId})" title="Hide comments">×</button>
                    </div>
                    <div class="comments-container" id="comments-container-${postId}">
                        <div class="loading-comments">Loading comments...</div>
                    </div>
                    <div class="add-comment-form">
                        <form onsubmit="submitComment(event, ${postId})">
                            <div class="comment-avatar" style="background-color: ${generateAvatarColor('current_user')}">
                                ${getCurrentUsername().charAt(0).toUpperCase()}
                            </div>
                            <div class="input-wrapper">
                                <input type="text" placeholder="Write a comment..." required maxlength="500">
                            </div>
                            <button type="submit" class="btn btn-primary btn-small">Post</button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    `;
}

// Infinite scroll setup
function setupInfiniteScroll() {
    const postsContainer = document.getElementById('posts');

    // Remove existing scroll listener to avoid duplicates
    window.removeEventListener('scroll', handleHomeScroll);

    // Add new scroll listener
    window.addEventListener('scroll', handleHomeScroll);
}

function handleHomeScroll() {
    // Only handle scroll if we're on the home section
    if (homeSection.style.display === 'none') return;

    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    const windowHeight = window.innerHeight;
    const documentHeight = document.documentElement.scrollHeight;

    // Load more when user is 100px from bottom
    if (scrollTop + windowHeight >= documentHeight - 100) {
        loadMoreHomePosts();
    }
}

async function loadMoreHomePosts() {
    if (!hasMoreHomePosts || isLoadingHomeFeed) return;

    currentHomePage++;
    console.log('Loading more posts, page:', currentHomePage);

    // Add loading indicator
    const postsContainer = document.getElementById('posts');
    const loadingDiv = document.createElement('div');
    loadingDiv.id = 'loading-more';
    loadingDiv.className = 'text-center';
    loadingDiv.innerHTML = '<p>Loading more posts...</p>';
    postsContainer.appendChild(loadingDiv);

    await loadHomeFeed(currentHomePage, true);

    // Remove loading indicator
    const loadingElement = document.getElementById('loading-more');
    if (loadingElement) {
        loadingElement.remove();
    }
}

// Comments functionality
async function toggleComments(postId) {
    console.log('Toggling comments for post:', postId); // Debug log
    const commentsSection = document.getElementById(`comments-${postId}`);
    const commentsContainer = document.getElementById(`comments-container-${postId}`);

    if (!commentsSection) {
        console.error('Comments section not found for post:', postId);
        return;
    }

    if (commentsSection.style.display === 'none') {
        commentsSection.style.display = 'block';

        // Load comments if not already loaded
        if (commentsContainer && commentsContainer.innerHTML.includes('Loading comments...')) {
            console.log('Loading comments for post:', postId); // Debug log
            await loadCommentsForPost(postId);
        }
    } else {
        commentsSection.style.display = 'none';
    }
}

async function loadCommentsForPost(postId) {
    const commentsContainer = document.getElementById(`comments-container-${postId}`);

    if (!commentsContainer) {
        console.error('Comments container not found for post:', postId);
        return;
    }

    try {
        console.log('Making API call to load comments for post:', postId); // Debug log
        const comments = await apiCall(`/post/comments?postId=${postId}`);
        console.log('Received comments:', comments); // Debug log

        if (comments && comments.length > 0) {
            commentsContainer.innerHTML = comments.map(comment => createCommentHTML(comment)).join('');

            // Update the comment count with total number of comments (including nested)
            const totalCommentCount = countTotalComments(comments);
            updateCommentCountDisplay(postId, totalCommentCount);
        } else {
            commentsContainer.innerHTML = '<div class="no-comments">No comments yet. Be the first to comment!</div>';
            updateCommentCountDisplay(postId, 0);
        }
    } catch (error) {
        console.error('Failed to load comments for post:', postId, error);
        commentsContainer.innerHTML = '<div class="error-comments">Failed to load comments. Please try again.</div>';
    }
}

function createCommentHTML(commentTree, postId = null, depth = 0) {
    // Handle CommentTree structure from API - use exact field names from backend
    const comment = commentTree.comment || commentTree.Comment || commentTree;
    const children = commentTree.children || commentTree.Children || [];

    console.log('Full comment object:', comment); // Debug log to see actual structure

    // Extract comment data using exact field names from CommentWithUser struct
    // Handle both the expected format (lowercase) and actual format (uppercase)
    const commentId = comment.id || comment.Id;
    const userId = comment.user_id || comment.UserId;
    const username = comment.username || comment.Username || `User${userId}`;
    const commentText = comment.message || comment.Message || '';
    const createdAt = comment.created_at || comment.CreatedAt || '';

    console.log('Extracted comment data:', {
        commentId,
        userId,
        username,
        commentText,
        createdAt
    }); // Debug log

    // Limit nesting depth to prevent UI overflow
    const maxDepth = 3;
    const currentDepth = Math.min(depth, maxDepth);
    const depthClass = currentDepth > 0 ? `nested-depth-${currentDepth}` : '';

    // Generate user avatar with better color scheme
    const avatarColor = generateAvatarColor(username);
    const timeAgo = getTimeAgo(createdAt);

    let childrenHTML = '';
    if (children && children.length > 0 && currentDepth < maxDepth) {
        childrenHTML = children.map(child => createCommentHTML(child, postId, depth + 1)).join('');
    }

    const replyFormId = `reply-form-${commentId}`;
    const escapedUsername = username.replace(/'/g, "\\'");
    const hasReplies = children && children.length > 0;

    return `
        <article class="comment-item ${depthClass}" id="comment-${commentId}" data-depth="${currentDepth}" data-comment-id="${commentId}">
            <div class="comment-main">
                <div class="comment-avatar-wrapper">
                    <div class="comment-avatar" style="background-color: ${avatarColor}">
                        ${username.charAt(0).toUpperCase()}
                    </div>
                    ${hasReplies && currentDepth === 0 ? '<div class="reply-thread-line"></div>' : ''}
                </div>

                <div class="comment-body">
                    <div class="comment-header">
                        <div class="comment-author-info">
                            <span class="comment-username">${username}</span>
                            <span class="comment-timestamp" title="${formatDate(createdAt)}">
                                ${timeAgo}
                            </span>
                        </div>
                        ${currentDepth < maxDepth ? `
                        <div class="comment-menu">
                            <button class="comment-menu-btn" onclick="toggleCommentMenu(${commentId})">
                                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                    <circle cx="12" cy="5" r="2"/>
                                    <circle cx="12" cy="12" r="2"/>
                                    <circle cx="12" cy="19" r="2"/>
                                </svg>
                            </button>
                        </div>
                        ` : ''}
                    </div>

                    <div class="comment-content">
                        <p class="comment-text">${commentText}</p>
                    </div>

                    ${currentDepth < maxDepth ? `
                    <div class="comment-actions">
                        <button class="comment-action-btn reply-btn" onclick="openReplyForm(${commentId}, '${escapedUsername}', ${postId || 'null'})">
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path>
                            </svg>
                            Reply
                        </button>
                        <button class="comment-action-btn like-btn" onclick="likeComment(${commentId})">
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
                            </svg>
                            <span class="like-count">0</span>
                        </button>
                        ${hasReplies ? `
                        <button class="comment-action-btn toggle-replies-btn" onclick="toggleReplies(${commentId})">
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <polyline points="6 9 12 15 18 9"></polyline>
                            </svg>
                            ${children.length} ${children.length === 1 ? 'reply' : 'replies'}
                        </button>
                        ` : ''}
                    </div>
                    ` : ''}
                </div>
            </div>

            ${currentDepth < maxDepth ? `
            <div class="comment-reply-form" id="${replyFormId}" style="display: none;">
                <div class="reply-form-container">
                    <div class="reply-form-header">
                        <span class="reply-indicator">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <path d="M3 20h18L12 4z"></path>
                            </svg>
                            Replying to <strong>@${username}</strong>
                        </span>
                        <button class="cancel-reply-btn" onclick="cancelReply(${commentId})">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <line x1="18" y1="6" x2="6" y2="18"></line>
                                <line x1="6" y1="6" x2="18" y2="18"></line>
                            </svg>
                        </button>
                    </div>
                    <form class="reply-form" onsubmit="submitReply(event, ${commentId}, ${postId || 'null'})">
                        <div class="reply-input-container">
                            <div class="reply-avatar" style="background-color: ${generateAvatarColor('current_user')}">
                                ${getCurrentUsername().charAt(0).toUpperCase()}
                            </div>
                            <div class="reply-input-wrapper">
                                <input type="text" placeholder="Write a thoughtful reply..." required maxlength="500">
                                <div class="reply-form-actions">
                                    <button type="submit" class="reply-submit-btn">
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                            <line x1="22" y1="2" x2="11" y2="13"></line>
                                            <polygon points="22,2 15,22 11,13 2,9 22,2"></polygon>
                                        </svg>
                                        Send
                                    </button>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            ` : ''}

            ${childrenHTML ? `
            <div class="nested-comments-container" id="replies-${commentId}">
                <div class="nested-comments">
                    ${childrenHTML}
                </div>
            </div>
            ` : ''}
        </article>
    `;
}

// Helper functions for enhanced comment UI
function generateAvatarColor(username) {
    const colors = [
        '#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FECA57',
        '#FF9FF3', '#54A0FF', '#5F27CD', '#00D2D3', '#FF9F43',
        '#10AC84', '#EE5A24', '#0984E3', '#6C5CE7', '#FD79A8'
    ];
    let hash = 0;
    for (let i = 0; i < username.length; i++) {
        hash = username.charCodeAt(i) + ((hash << 5) - hash);
    }
    return colors[Math.abs(hash) % colors.length];
}

function getTimeAgo(dateString) {
    if (!dateString) return 'now';

    const now = new Date();
    const date = new Date(dateString);
    const diffInSeconds = Math.floor((now - date) / 1000);

    if (diffInSeconds < 60) return 'now';
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m`;
    if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h`;
    if (diffInSeconds < 604800) return `${Math.floor(diffInSeconds / 86400)}d`;

    return formatDate(dateString);
}

function getCurrentUsername() {
    // Get the current user's username from localStorage or context
    return localStorage.getItem('username') || 'You';
}

function toggleCommentMenu(commentId) {
    // Toggle comment menu dropdown (placeholder for future features)
    console.log('Toggle menu for comment:', commentId);
}

function likeComment(commentId) {
    // Like/unlike comment functionality (placeholder)
    console.log('Like comment:', commentId);

    // Find the like button and update count
    const likeBtn = document.querySelector(`#comment-${commentId} .like-btn`);
    if (likeBtn) {
        const countSpan = likeBtn.querySelector('.like-count');
        if (countSpan) {
            const currentCount = parseInt(countSpan.textContent) || 0;
            countSpan.textContent = currentCount + 1;
            likeBtn.classList.add('liked');
        }
    }
}

function toggleReplies(commentId) {
    const repliesContainer = document.getElementById(`replies-${commentId}`);
    const toggleBtn = document.querySelector(`#comment-${commentId} .toggle-replies-btn`);

    if (repliesContainer && toggleBtn) {
        const isHidden = repliesContainer.style.display === 'none';
        repliesContainer.style.display = isHidden ? 'block' : 'none';

        const icon = toggleBtn.querySelector('svg');
        if (icon) {
            icon.style.transform = isHidden ? 'rotate(180deg)' : 'rotate(0deg)';
        }
    }
}

async function submitComment(event, postId) {
    event.preventDefault();
    console.log('Submitting comment for post:', postId); // Debug log

    const form = event.target;
    const input = form.querySelector('input');
    const comment = input.value.trim();

    if (!comment) {
        console.log('Empty comment, not submitting'); // Debug log
        return;
    }

    console.log('Comment text:', comment); // Debug log

    try {
        // Use exact field names from CommentRequest DTO
        const response = await apiCall('/addComment/post', {
            method: 'POST',
            body: JSON.stringify({
                parent_id: postId,
                parent_type: 'post',
                comment: comment
            })
        });

        console.log('Comment submission response:', response); // Debug log

        // Clear the input
        input.value = '';

        // Reload comments for this post (this will automatically update the comment count)
        await loadCommentsForPost(postId);

        showToast('Comment added successfully!');
    } catch (error) {
        console.error('Failed to add comment:', error);
        showToast('Failed to add comment. Please try again.', 'error');
    }
}

// Reply to comment functions
function openReplyForm(commentId, username, postId) {
    // Hide all other reply forms first
    document.querySelectorAll('.comment-reply-form').forEach(form => {
        form.style.display = 'none';
    });

    // Show the reply form for this comment
    const replyForm = document.getElementById(`reply-form-${commentId}`);
    if (replyForm) {
        replyForm.style.display = 'block';
        const input = replyForm.querySelector('input');
        if (input) {
            input.focus();
        }
    }
}

function cancelReply(commentId) {
    const replyForm = document.getElementById(`reply-form-${commentId}`);
    if (replyForm) {
        replyForm.style.display = 'none';
        const input = replyForm.querySelector('input');
        if (input) {
            input.value = '';
        }
    }
}

async function submitReply(event, parentCommentId, postId) {
    event.preventDefault();
    const form = event.target;
    const input = form.querySelector('input');
    const comment = input.value.trim();

    if (!comment) return;

    try {
        // Use exact field names from CommentRequest DTO
        await apiCall('/addComment/comment', {
            method: 'POST',
            body: JSON.stringify({
                parent_id: parentCommentId,
                parent_type: 'comment',
                comment: comment
            })
        });

        // Clear the input and hide the form
        input.value = '';
        cancelReply(parentCommentId);

        // Reload comments for this post to show the new reply (this will automatically update the comment count)
        if (postId) {
            await loadCommentsForPost(postId);
        }

        showToast('Reply added successfully!');
    } catch (error) {
        console.error('Failed to add reply:', error);
        showToast('Failed to add reply. Please try again.', 'error');
    }
}

// Helper function to count total comments including nested ones
function countTotalComments(commentTrees) {
    if (!commentTrees || !Array.isArray(commentTrees)) {
        return 0;
    }

    let totalCount = 0;

    for (const commentTree of commentTrees) {
        // Count the current comment
        totalCount += 1;

        // Recursively count children
        const children = commentTree.children || commentTree.Children || [];
        if (children.length > 0) {
            totalCount += countTotalComments(children);
        }
    }

    return totalCount;
}

function updateCommentCountDisplay(postId, totalCount) {
    const commentCountElement = document.querySelector(`#post-${postId} .comment-count`);
    if (commentCountElement) {
        commentCountElement.textContent = totalCount;
    }
}

function updateCommentCount(postId) {
    const commentCountElement = document.querySelector(`#post-${postId} .comment-count`);
    if (commentCountElement) {
        const currentCount = parseInt(commentCountElement.textContent) || 0;
        commentCountElement.textContent = currentCount + 1;
    }
}

async function loadMyPosts(page = 1, limit = 10) {
    try {
        console.log('Loading my posts...', { page, limit }); // Debug log
        const data = await apiCall(`/myPosts?page=${page}&limit=${limit}`);
        console.log('My posts data received:', data); // Debug log

        const myPostsList = document.getElementById('myPostsList');
        const postsStats = document.getElementById('postsStats');
        const pagination = document.getElementById('myPostsPagination');

        // Update stats using exact field names from HomepageResponse DTO
        const totalPosts = data.total_posts || 0;
        postsStats.innerHTML = `
            <div class="stats-item">
                <span class="stats-number">${totalPosts}</span>
                <span class="stats-label">Total Posts</span>
            </div>
        `;

        if (data.posts && data.posts.length > 0) {
            console.log('First post data:', data.posts[0]); // Debug: see what fields are available
            myPostsList.innerHTML = data.posts.map(post => {
                // Use exact field names from PostWithAuthor DTO
                const authorUsername = post.author_username;
                const postId = post.id;
                const postUrl = post.url;
                const postCaption = post.caption;
                const createdAt = post.created_at;
                const likeCount = post.like_count;
                const commentCount = post.comment_count;

                console.log('Processing post:', { postId, authorUsername, postUrl }); // Debug log

                return `
                <div class="post-card fade-in">
                    <div class="post-header">
                        <div class="post-author">
                            <div class="author-avatar">${authorUsername.charAt(0).toUpperCase()}</div>
                            <div class="author-info">
                                <div class="author-name">${authorUsername}</div>
                                <div class="post-date">${formatDate(createdAt)}</div>
                            </div>
                        </div>
                        <div class="post-actions">
                            <button class="btn btn-danger btn-small" onclick="deletePost(${postId})">
                                <i class="fas fa-trash"></i> Delete
                            </button>
                        </div>
                    </div>
                    <div class="post-image">
                        <img src="${postUrl}" alt="Post image" onerror="this.src='https://via.placeholder.com/400x400?text=Image+Not+Found'">
                    </div>
                    <div class="post-content">
                        <div class="post-caption">${postCaption}</div>
                        <div class="post-stats">
                            <button class="action-btn" onclick="likePost(${postId})">
                                <i class="fas fa-heart"></i> ${likeCount}
                            </button>
                            <button class="action-btn" onclick="openCommentModal(${postId}, 'post')">
                                <i class="fas fa-comment"></i> ${commentCount}
                            </button>
                        </div>
                    </div>
                </div>
            `;
            }).join('');

            // Show pagination if there are multiple pages using exact field names
            const totalPages = data.total_pages || 1;
            const currentPage = data.page || page;
            const hasNext = data.has_next || false;
            const hasPrev = data.has_prev || false;

            if (totalPages > 1) {
                pagination.style.display = 'block';
                pagination.innerHTML = `
                    <button class="btn btn-secondary"
                            onclick="loadMyPosts(${currentPage - 1})"
                            ${!hasPrev ? 'disabled' : ''}>
                        Previous
                    </button>
                    <span>Page ${currentPage} of ${totalPages}</span>
                    <button class="btn btn-secondary"
                            onclick="loadMyPosts(${currentPage + 1})"
                            ${!hasNext ? 'disabled' : ''}>
                        Next
                    </button>
                `;
            } else {
                pagination.style.display = 'none';
            }
        } else {
            myPostsList.innerHTML = '<div class="text-center"><p>You haven\'t created any posts yet. Create your first post!</p></div>';
            pagination.style.display = 'none';
        }
    } catch (error) {
        console.error('Failed to load my posts:', error);
        const myPostsList = document.getElementById('myPostsList');
        const postsStats = document.getElementById('postsStats');

        // Reset stats on error
        postsStats.innerHTML = `
            <div class="stats-item">
                <span class="stats-number">0</span>
                <span class="stats-label">Total Posts</span>
            </div>
        `;

        myPostsList.innerHTML = '<div class="text-center"><p>Failed to load posts. Please try again.</p></div>';
        showToast('Failed to load posts. Please try again.', 'error');
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

        console.log('Followers data received:', data); // Debug log

        if (data.followers && data.followers.length > 0) {
            console.log('Rendering followers with unfollow buttons:', data.followers); // Debug log
            followersList.innerHTML = data.followers.map(username => `
                <div class="follower-item fade-in">
                    <div class="follower-avatar">${username.charAt(0).toUpperCase()}</div>
                    <div class="follower-info">
                        <div class="follower-name">${username}</div>
                        <button class="btn btn-danger btn-small" onclick="openUnfollowModal('${username}')">
                            Unfollow
                        </button>
                    </div>
                </div>
            `).join('');
        } else {
            console.log('No followers found, showing empty message'); // Debug log
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

        // Update like count in UI
        updateLikeCount(postId);
        showToast('Post liked!');
    } catch (error) {
        console.error('Failed to like post:', error);
        showToast('Failed to like post. Please try again.', 'error');
    }
}

function updateLikeCount(postId) {
    const likeCountElement = document.querySelector(`#post-${postId} .like-count`);
    if (likeCountElement) {
        const currentCount = parseInt(likeCountElement.textContent) || 0;
        likeCountElement.textContent = currentCount + 1;

        // Add visual feedback
        const likeBtn = likeCountElement.closest('.like-btn');
        if (likeBtn) {
            likeBtn.classList.add('liked');
            setTimeout(() => likeBtn.classList.remove('liked'), 1000);
        }
    }
}

async function deletePost(postId) {
    if (!confirm('Are you sure you want to delete this post? This action cannot be undone.')) {
        return;
    }

    try {
        await apiCall(`/deletePost/${postId}`, {
            method: 'DELETE'
        });

        showToast('Post deleted successfully!');
        // Reload my posts if we're on that section
        if (myPostsSection.style.display !== 'none') {
            loadMyPosts();
        }
    } catch (error) {
        console.error('Failed to delete post:', error);
        showToast('Failed to delete post. Please try again.', 'error');
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

async function unfollowUser(username) {
    try {
        await apiCall('/unfollow', {
            method: 'POST',
            body: JSON.stringify({ username })
        });

        showToast(`You unfollowed ${username}`);
        closeModal(unfollowModal);
        if (followersSection.style.display !== 'none') {
            loadFollowers();
        }
    } catch (error) {
        console.error('Failed to unfollow user:', error);
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

function openUnfollowModal(username = '') {
    document.getElementById('unfollowUsername').value = username;
    document.getElementById('unfollowConfirmText').textContent = `Are you sure you want to unfollow ${username}?`;
    openModal(unfollowModal);
}

function openCommentModal(parentId, parentType) {
    document.getElementById('commentParentId').value = parentId;
    document.getElementById('commentParentType').value = parentType;
    document.getElementById('commentText').value = '';
    openModal(commentModal);
}

// Event Listeners
document.addEventListener('DOMContentLoaded', function () {
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
    myPostsBtn.addEventListener('click', showMyPosts);
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

    // Unfollow form
    document.getElementById('unfollowForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('unfollowUsername').value;
        await unfollowUser(username);
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
window.openUnfollowModal = openUnfollowModal;
window.openCommentModal = openCommentModal;
window.likePost = likePost;
window.deletePost = deletePost;
window.loadMyPosts = loadMyPosts;
window.closeModal = closeModal;
window.toggleComments = toggleComments;
window.submitComment = submitComment;
window.loadMyPosts = loadMyPosts;
window.closeModal = closeModal;

// Debug function to check token status
window.debugTokens = function () {
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
window.testUsersAPI = async function () {
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
