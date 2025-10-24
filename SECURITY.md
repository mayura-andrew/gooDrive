# Security Policy

Security guidelines for gooDrive users and contributors.

---

## Reporting Security Vulnerabilities

**DO NOT** open public GitHub issues for security vulnerabilities.

### Report Privately

Please report security vulnerabilities to:

**Email**: security@mayura-andrew.dev

**Include in Report:**
1. Description of vulnerability
2. Steps to reproduce
3. Potential impact
4. Suggested fix (if any)

**Response Time:**
- Critical: Within 24 hours
- High: Within 3 days
- Medium: Within 7 days
- Low: Within 14 days

---

## Security Features

### OAuth 2.0 Authentication

✅ **Secure Implementation:**
- Uses OAuth 2.0 authorization code flow
- No password stored locally
- Refresh tokens for long-term access
- Automatic token rotation

⚠️ **Token Storage:**
- Tokens stored in `~/.config/gooDrive/.tokens.json`
- File permissions: `0600` (read/write for owner only)
- Never stored in version control
- Encrypted if using system keychain (future feature)

### HTTPS Only

✅ **All Communications:**
- OAuth flow uses HTTPS
- Google Drive API calls use HTTPS
- Local callback server uses HTTP (for localhost only)
- No data transmitted in plaintext

### No Credentials in Code

✅ **Safe Practices:**
- No hardcoded API keys
- No embedded credentials
- OAuth.json not committed to Git
- Credentials loaded from environment or file

---

## Best Practices for Users

### 1. OAuth Credentials

```bash
# ✅ DO: Place oauth.json outside project
~/.config/gooDrive/oauth.json

# ❌ DON'T: Commit to Git
git add oauth.json  # DON'T DO THIS!

# ✅ DO: Add to .gitignore
echo "oauth.json" >> .gitignore
echo ".drive-cli-token.json" >> .gitignore
```

### 2. Token Management

```bash
# ✅ DO: Use automatic token refresh
# Tokens refresh automatically on use

# ✅ DO: Periodically re-authenticate
gooDrive auth refresh

# ❌ DON'T: Manually edit token file
# Use 'gooDrive auth' instead

# ✅ DO: Keep tokens private
# Tokens stored with 0600 permissions automatically
```

### 3. File Sharing

```bash
# ✅ DO: Share with specific users
gooDrive share <file-id> --email user@example.com --role reader

# ✅ DO: Use restricted sharing roles
# reader    - View only
# commenter - View and comment
# writer    - Edit

# ❌ DON'T: Make files publicly accessible
# Without proper permissions review
```

### 4. Sync Operations

```bash
# ✅ DO: Test before syncing large directories
gooDrive list

# ✅ DO: Use selective sync
# Only sync needed folders

# ✅ DO: Backup important files
cp -r ~/important_folder ~/important_folder.backup

# ❌ DON'T: Sync system directories
# Don't sync /usr, /etc, or system files
```

---

## Security Checklist

### Installation
- [ ] Download from official GitHub releases
- [ ] Verify checksums match official file
- [ ] Don't run from untrusted sources
- [ ] Keep Go updated to latest version

### Setup
- [ ] Create oauth.json in secure location
- [ ] Set file permissions to 0600
- [ ] Don't share oauth.json with others
- [ ] Use strong Google account password

### Usage
- [ ] Keep gooDrive updated
- [ ] Monitor file sharing permissions
- [ ] Review audit logs periodically
- [ ] Don't share your config directory

### Storage
- [ ] Tokens stored securely (0600)
- [ ] Config files not world-readable
- [ ] Don't backup tokens to cloud without encryption
- [ ] Use encrypted storage for sensitive projects

---

## Known Security Limitations

### Current Version (v0.1.0-beta.1)

⚠️ **Limitations:**
1. No end-to-end encryption for files
2. Token stored in plaintext (but restricted permissions)
3. No built-in audit logging
4. Single-device token management only
5. No IP-based access restrictions

### Planned Improvements

✅ **Future Versions:**
- System keychain integration for token storage
- Optional end-to-end encryption
- Comprehensive audit logging
- Multi-device token management
- Rate limiting and access controls
- Signed binary releases with GPG

---

## Dependency Security

### Dependency Updates

```bash
# Check for vulnerabilities
go list -json -m all | nancy sleuth

# Update dependencies
go get -u ./...

# Check for breaking changes
go mod tidy
go mod verify
```

### Third-Party Libraries

**Current Dependencies:**
- `golang.org/x/oauth2` - OAuth 2.0 support
- `google.golang.org/api` - Google Drive API
- `github.com/spf13/cobra` - CLI framework
- `github.com/fsnotify/fsnotify` - File watching

All libraries are from official sources and regularly updated.

---

## Development Security

### For Contributors

1. **Code Review**
   - All PRs require security review
   - No direct commits to main
   - Two-person approval required

2. **Secrets Management**
   - Never commit secrets
   - Use environment variables in CI
   - Rotate tokens regularly

3. **Dependency Management**
   - Only use trusted dependencies
   - Verify source before including
   - Keep dependencies updated

### Security Testing

```bash
# Run security checks
go vet ./...
golangci-lint run
gosec ./...

# Check for vulnerabilities
go list -json -m all | nancy sleuth
```

---

## Incident Response

### If You Discover a Vulnerability

1. **Stop using affected version** (if critical)
2. **Report privately** (see top of this document)
3. **Wait for patch** (don't disclose publicly)
4. **Update when released**
5. **Re-authenticate** if token compromise suspected

### If Credentials Are Compromised

```bash
# 1. Revoke access immediately
# Go to: https://myaccount.google.com/security-checkup
# Remove gooDrive from authorized apps

# 2. Remove token file
rm ~/.config/gooDrive/.tokens.json

# 3. Re-authenticate with fresh credentials
gooDrive auth

# 4. Review Google Drive activity
# Check for unauthorized changes
```

---

## Security Transparency

### Source Code Audit

The source code is publicly available for review:
- https://github.com/mayura-andrew/gooDrive
- All code changes tracked in Git history
- Security improvements documented in CHANGELOG

### Binary Verification

```bash
# Verify binary checksums
sha256sum -c checksums.txt

# Verify GPG signature (when available)
gpg --verify gooDrive.sig gooDrive
```

### Responsible Disclosure

We follow responsible disclosure practices:
- Private initial report
- Time to patch before disclosure
- Public credit for reporter (if desired)
- Detailed security advisory published

---

## Compliance

### Data Privacy

✅ **What gooDrive Respects:**
- Google Privacy Policy
- GDPR requirements
- Your data ownership
- Your privacy settings

✅ **What We Don't Do:**
- Collect usage statistics
- Track your files
- Send data to third parties
- Store files on our servers

### Google Drive API

✅ **OAuth Scopes:**
- Only requests `https://www.googleapis.com/auth/drive`
- Standard scope for full Drive access
- Follows Google's API policies
- Respects user permissions

---

## Updates & Patches

### Security Updates

Subscribe to security updates:
- Watch GitHub repository: https://github.com/mayura-andrew/gooDrive
- Enable notifications for releases
- Follow changelog for security fixes

### Installation Updates

```bash
# Update to latest
go install github.com/mayura-andrew/goodrive@latest

# Update to specific security version
go install github.com/mayura-andrew/goodrive@v0.1.1
```

---

## FAQ

### Q: Is my password stored by gooDrive?
**A:** No. gooDrive uses OAuth, you sign in directly with Google. Your password never reaches gooDrive.

### Q: Where are tokens stored?
**A:** In `~/.config/gooDrive/.tokens.json` with restricted permissions (0600).

### Q: Can I use gooDrive from multiple devices?
**A:** Each device gets its own token. Each device must authenticate separately.

### Q: What if I lose my device?
**A:** Revoke access through Google account settings. The token on the lost device becomes useless.

### Q: Is this tool officially by Google?
**A:** No. It's a community tool for Google Drive. Not affiliated with Google.

### Q: Should I trust gooDrive?
**A:** It's open source, so you can audit the code. Use at your own risk and keep backups.

---

## Additional Resources

- [Google OAuth 2.0 Security](https://developers.google.com/identity/protocols/oauth2/security)
- [OWASP Security Guidelines](https://owasp.org/)
- [Go Security Best Practices](https://go.dev/doc/security)
- [GitHub Security](https://github.blog/security/)

---

**Last Updated:** October 24, 2025
**Version:** 1.0
