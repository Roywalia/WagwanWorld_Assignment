# Task: Build a Public RSVP Page

## Your Mission

Design and implement a public-facing RSVP page where guests can register for an event.

## What to Build

### Public RSVP Page (`/rsvp` route)

Create a new page that allows anyone to RSVP for the event.

**Required Features:**

1. **Event Information Display**
   - Event title
   - Date and time
   - Location or description
   - Make it visually appealing and UX friendly

2. **RSVP Form**
   - Name
   - Email
   - Phone
   - RSVP Status selection:
     - "Yes, I'll be there!" (attending)
     - "Maybe" (pending)
     - "Sorry, can't make it" (declined)

3. **Form Behavior**
   - Validate all inputs before submission
   - Submit to the existing API
   - Show success message after submission
   - Show friendly error messages if something goes wrong
   - Clear form after successful submission

4. **Design Requirements**
   - Mobile-responsive
   - Use TailwindCSS
   - Different look from the admin interface
   - Professional and polished

## Important Notes

⚠️ **The existing codebase is functional but not perfect.**

While implementing your RSVP page, you may encounter issues with the existing code. Part of this task is to:
- Identify any bugs or problems
- Fix them as you go
- Document what you fixed and why

## Database Migration Task

Implement a database migration system to handle schema changes.

### Your Mission
Choose and implement a database migration system for the application. You have the freedom to select any migration tool or framework that you're comfortable with.

### Migration Requirements

You need to implement the following database schema changes:

1. **Add Guest Notes Field**
   - Add a `notes` column to the `guests` table
   - Column should be `TEXT` type and nullable
   - Update the Go struct and queries to handle the new field
   - Update the API to accept and return notes in guest responses

2. **Add Events Table**
   - Create a new `events` table with:
     - `id` (SERIAL PRIMARY KEY)
     - `title` (VARCHAR(255) NOT NULL)
     - `description` (TEXT)
     - `event_date` (TIMESTAMP)
     - `location` (VARCHAR(255))
     - `created_at` (TIMESTAMP)
   - Add a foreign key `event_id` to the `guests` table
   - Update the guest creation to require an event_id

3. **Add RSVP Tracking**
   - Add `rsvp_date` column to guests table
   - Add `plus_ones` integer column (default 0)
   - Add `dietary_restrictions` text column
   - Update the API to handle these new fields

### Migration Tool Options

Choose one of these approaches (or suggest your own):

#### Option 1: golang-migrate
- **Tool**: [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- **Pros**: Popular, well-documented, supports many databases
- **Usage**: `migrate create -ext sql -dir migrations -seq add_guest_notes`
- **Commands**: `migrate -path migrations -database "postgres://..." up`

#### Option 2: Atlas
- **Tool**: [ariga/atlas](https://atlasgo.io/)
- **Pros**: Modern, declarative, great for complex schemas
- **Usage**: `atlas migrate diff add_guest_notes --env local`
- **Commands**: `atlas migrate apply --env local`

#### Option 3: SQLC with Custom Migrations
- **Tool**: [sqlc](https://sqlc.dev/) + custom migration runner
- **Pros**: Type-safe SQL, generates Go code
- **Usage**: Write SQL files, use sqlc to generate Go code

#### Option 4: Custom Migration System
- **Tool**: Build your own migration runner
- **Pros**: Full control, learning experience
- **Requirements**: Track applied migrations, support up/down, error handling

#### Option 5: Prisma (if you want to try something different)
- **Tool**: [Prisma](https://www.prisma.io/) with Go
- **Pros**: Type-safe, great developer experience
- **Usage**: `prisma migrate dev --name add_guest_notes`

#### Option 6: Flyway (Java-based, but can be used)
- **Tool**: [Flyway](https://flywaydb.org/)
- **Pros**: Enterprise-grade, supports many databases
- **Usage**: `flyway migrate`

### Tool Selection Tips

- **For beginners**: Start with golang-migrate (Option 1)
- **For modern approaches**: Try Atlas (Option 2) 
- **For type safety**: Consider SQLC (Option 3)
- **For learning**: Build custom system (Option 4)
- **For enterprise**: Consider Flyway (Option 6)

### Implementation Requirements

1. **Migration Files**
   - Create migration files for each schema change
   - Include both forward and rollback migrations
   - Use consistent naming convention

2. **Migration Runner**
   - Apply migrations in correct order
   - Track which migrations have been applied
   - Support rollback functionality
   - Handle errors gracefully

3. **Integration**
   - Update Go structs to match new schema
   - Update database queries
   - Update API endpoints
   - Ensure backward compatibility

4. **Documentation**
   - Document your chosen approach
   - Explain how to run migrations
   - Include rollback instructions
   - Document any challenges faced

### Testing Requirements

- Test that migrations apply in correct order
- Test rollback functionality
- Ensure existing data is preserved
- Test that the application works after migrations
- Test error handling (what happens if a migration fails?)

### Additional Challenges (Optional)

If you want to go further, consider these challenges:

1. **Data Migration**
   - How do you handle existing guests when adding the `event_id` foreign key?
   - What if you need to migrate data from one column to another?

2. **Zero-Downtime Deployments**
   - How do you deploy schema changes without breaking the application?
   - Consider blue-green deployment strategies

3. **Schema Validation**
   - How do you ensure the database schema matches your expectations?
   - Add validation to check schema consistency

4. **Migration Dependencies**
   - What if one migration depends on another?
   - How do you handle complex migration chains?

5. **Performance Considerations**
   - How do you handle migrations on large tables?
   - Consider adding indexes during migrations

### Evaluation Criteria

We'll evaluate based on:
- **Tool Selection**: Did you choose an appropriate tool for the job?
- **Implementation**: Is the migration system robust and well-implemented?
- **Code Quality**: Is the code clean and maintainable?
- **Documentation**: Is the system well-documented?
- **Testing**: Did you thoroughly test the migration system?
- **Problem Solving**: How did you handle challenges and edge cases?

## Bonus Points (Optional)

If you have extra time and want to go further:

- Add duplicate email detection (prevent same email registering twice)
- Add a "View Your RSVP" feature (lookup by email)
- Add form animations or transitions
- Display attendee count ("X people attending")
- Add dietary restrictions or notes field
- Implement better error handling

## Testing Your Work

Make sure to test:
- ✅ Form validation works (required fields, email format)
- ✅ Successful submission saves to database
- ✅ Success message displays correctly
- ✅ Error messages display if API fails
- ✅ Works on mobile (responsive design)
- ✅ Admin page can see newly created guests
- ✅ All existing functionality still works

### Migration Testing Requirements

- ✅ Migration system initializes correctly
- ✅ All three migrations apply in order without errors
- ✅ Database schema changes are applied correctly
- ✅ Rollback functionality works (migrate down)
- ✅ Migration status shows correct information
- ✅ Application works with new schema after migrations
- ✅ Existing data is preserved during migrations
- ✅ New API endpoints handle additional fields

## Submission

When complete, submit:

1. **Your code** (zip file or private GitHub repo)

2. **SUBMISSION.md** with:
   - Bugs you found and fixed
   - Features you implemented
   - Challenges you faced
   - How to test your work
   - What you'd improve with more time
   - **Migration system details**:
     - Which migration tool you chose and why
     - How to set up and run migrations
     - Database schema changes made
     - Any challenges with data migration
     - How you handled rollbacks
     - Performance considerations (if any)

## Evaluation

We're looking for:
- **Working code** - Does it work as specified?
- **Code quality** - Is it clean, readable, and maintainable?
- **Problem-solving** - Did you identify and fix issues?
- **Design** - Is the RSVP page well-designed and user-friendly?
- **Testing** - Does everything work correctly?
- **Migration system** - Is the database migration system robust and well-implemented?
- **Database design** - Are the schema changes logical and well-structured?
- **Documentation** - Is the migration system well-documented and easy to use?

## Time Estimate

This should take approximately 24-48 hours depending on your experience level and migration tool choice.

**Time breakdown:**
- RSVP page implementation: 12-20 hours
- Database migration system: 8-16 hours (varies by tool choice)
- Testing and debugging: 2-6 hours
- Documentation: 2-6 hours

**Migration tool time estimates:**
- **golang-migrate**: 8-12 hours (well-documented, straightforward)
- **Atlas**: 10-14 hours (modern but requires learning)
- **SQLC**: 12-16 hours (more complex setup)
- **Custom system**: 16-24 hours (full control but more work)
- **Prisma**: 10-14 hours (if familiar with Prisma)
- **Flyway**: 8-12 hours (enterprise tool, good docs)

Take your time, write clean code, and don't hesitate to ask questions if you get stuck!

Good luck! 🚀

