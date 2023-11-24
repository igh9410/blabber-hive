import React, { ChangeEvent, useState } from 'react';
import styles from './SignUpForm.module.scss';

export function SignUpForm() {
  const [imagePreviewUrl, setImagePreviewUrl] = useState('');

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
      <form id={styles.signupForm}>
        <h2>Sign Up</h2>
        <div className={styles.formGroup}>
          <label htmlFor="username">Username (required)</label>
          <input />
        </div>
        <div className="formGroup">
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
      </form>
    </div>
  );
}
