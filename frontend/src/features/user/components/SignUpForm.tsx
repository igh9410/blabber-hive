import React, { ChangeEvent, useRef, useState } from 'react';
import styles from './SignUpForm.module.scss';
import { Form } from 'react-router-dom';
import { useSignUp } from '@hooks';

export function SignUpForm() {
  const { signUp } = useSignUp();
  const [imagePreviewUrl, setImagePreviewUrl] = useState('');
  const [usernameError, setUsernameError] = useState(''); // State for username error message
  const usernameRef = useRef<HTMLInputElement | null>(null);

  const handleFormSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // You can now use the username state directly
    // Accessing the username value
    const username = usernameRef.current?.value;

    if (!username) {
      setUsernameError('Username is required'); // Set error message
      return;
    }

    // Reset error message if username is provided
    setUsernameError('');
    signUp({ username });
  };

  const handleImageChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const reader = new FileReader();
      reader.onload = (event: ProgressEvent<FileReader>) => {
        setImagePreviewUrl(event.target?.result as string);
      };
      reader.readAsDataURL(e.target.files[0]);
    }
  };
  return (
    <div className={styles.signupFormContainer}>
      <Form method="POST" id={styles.signupForm} onSubmit={handleFormSubmit}>
        <h2>Sign Up</h2>
        <div className={styles.formGroup}>
          <label htmlFor="username">Username (required)</label>
          <input type="text" id="username" name="username" ref={usernameRef} />
          {usernameError && (
            <div className={styles.error}>{usernameError}</div>
          )}{' '}
          {/* Display error message */}
        </div>
        <div className={styles.formGroup}>
          <label htmlFor="profileImage">Profile Image (optional)</label>
          <input
            type="file"
            id="profileImage"
            name="profileImage"
            accept="image/*"
            onChange={handleImageChange}
          />
        </div>
        <div className={styles.formGroup}>
          {imagePreviewUrl && (
            <img
              id="imagePreview"
              src={imagePreviewUrl}
              alt="Image Preview"
              style={{ display: 'block', maxWidth: '100%', height: 'auto' }}
            />
          )}
        </div>
        <button type="submit">Sign Up</button>
      </Form>
    </div>
  );
}
