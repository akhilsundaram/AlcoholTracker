<script>
  import { onMount } from "svelte";
  import Modal from './Modal.svelte';
  import { AddTrackerEntry, GetEntriesByDate, GetAlcoholCategories,ValidateFormDate, GetDrinks, GetDaysSinceLastDrink, GetDrinkTagColor, GetDrinkCount } from "../wailsjs/go/main/App";

  let year = new Date().getFullYear();
  let month = new Date().getMonth() + 1;
  let day = new Date().getDate();
  let category = "";
  let quantity = 0;
  let cost = 0;
  let entries = [];
  let activeTab = "calendar";
  let alcoholCategories = [];
  let validDate = false
  let daysSinceLastDrink = 0;
  let drinksToday = 0.0;
  let progress = 0; // This should be between 0 and 100
  let tag = { text: "N/A", color: "gray" };
  const tagColors = ["gray", "#60aa9b", "#43766c", "#ffdf60", "#fa8072", "#ed4d09"]
  
  $: daysDisplay = (daysSinceLastDrink == -1) ? "∞" : daysSinceLastDrink;
  /* Calendar STUFF */
  let currentYear = new Date().getFullYear();

  const months = [
    "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul",
    "Aug", "Sep", "Oct", "Nov", "Dec"
  ];

  function isLeapYear(year) {
    return (year % 4 === 0 && year % 100 !== 0) || (year % 400 === 0);
  }

  function getDaysInMonth(year) {
    return {
      Jan: 31, 
      Feb: isLeapYear(year) ? 29 : 28, 
      Mar: 31, Apr: 30, May: 31, Jun: 30, Jul: 31,
      Aug: 31, Sep: 30, Oct: 31, Nov: 30, Dec: 31
    };
  }

  let daysInMonth = getDaysInMonth(currentYear);

  function changeYear(direction) {
    currentYear += direction;
    daysInMonth = getDaysInMonth(currentYear);
    return false;
  }

  function handleCalendarButton(event) {
    if (event.type === "click" || event.key === "Enter" || event.code === "Space") {
      const year = Number(event.currentTarget.dataset.year);
      const month = Number(event.currentTarget.dataset.month);
      const day = Number(event.currentTarget.dataset.day);
      
      console.log(`Clicked/Pressed on ${year}-${Number(month) + 1}-${day}`);
      GetDrinks(Number(year),Number(month),Number(day))
      let x = GetDaysSinceLastDrink()
      // console.log(`${x}`);
      openModal(year,month,day)
    }
  }

  async function calendarDrinksMap(year, month, day) {
    let drink = await GetDrinks(Number(year),Number(month),Number(day))
    return `drinks-${drink}`
  }

  /* END CALENDAR */


  // MODAL
  let showModal = false;
  let selectedYear = 2024;
  let selectedMonth = 1;
  let selectedDay = 1;

  function openModal(year, month, day) {
    showModal = true;
    selectedYear = year;
    selectedMonth = month;
    selectedDay = day;

  }

  function closeModal() {
    showModal = false;
  }

  function onModalCloseCallback() {
    Refresh();
  }


  // END MODAL

  async function addEntry() {
    validDate = await ValidateFormDate(year,month,day)
    if (!category || quantity <= 0 || cost <= 0) {
      alert("Please enter valid details.");
      return;
    }
    if (!validDate) {
      alert("Invalid Date");
      return;
    }

    try {
      await AddTrackerEntry(year, month, day, category, quantity, cost);
      await fetchEntries();
      await Refresh();
    } catch (err) {
      console.error("Error adding entry:", err);
    }
  }

  async function Refresh(){
    changeYear(-1);
    changeYear(+1);
    daysSinceLastDrink = await GetDaysSinceLastDrink();
    const today = new Date();
    const year = today.getFullYear();
    const month = today.getMonth() + 1; // Months are zero-based (0 = January)
    const day = today.getDate();
    let drinkBar = 0;

    drinksToday = await GetDrinkCount(year,month,day)
    drinkBar = await GetDrinkTagColor(year,month,day);
    let tagText = await GetDrinks(year,month,day);
    tag = { text: tagText, color: tagColors[drinkBar] };
    progress = Math.min((drinkBar / 5) * 100, 100);
  }


  async function getAlcoholTypes() {
    try {
      alcoholCategories = await GetAlcoholCategories();
    } catch (err) {
      console.error("Error fetching alcohol categories:", err);
    }
  }

  async function fetchEntries() {
    try {
      entries = await GetEntriesByDate(year, month, day, category);
    } catch (err) {
      console.error("Error fetching entries:", err);
    }
  }


  function setTab(tab) {
    activeTab = tab;
  }



  // ON MOUNT
  onMount(getAlcoholTypes);

  // Scroll bar hidden but scrollable
  onMount(() => {
    document.documentElement.style.overflow = 'auto'; // Enable scrolling
    document.body.style.overflow = 'auto';
    document.body.style.margin = '0';
    document.body.style.padding = '0';

    // Hide scrollbar for Webkit browsers (Chrome, Safari, Edge)
    document.documentElement.style.scrollbarWidth = 'none'; // Firefox
    document.documentElement.style.setProperty("scrollbar-width", "none");
    document.documentElement.style.setProperty("scroll-behavior", "smooth");

    document.body.style.setProperty("-ms-overflow-style", "none"); // IE & Edge
    document.body.style.setProperty("scrollbar-width", "none"); // Firefox

    let style = document.createElement('style');
    style.innerHTML = `
      ::-webkit-scrollbar {
        display: none !important;
      }
    `;
    document.head.appendChild(style);
  });

  onMount(async () => {
    daysSinceLastDrink = await GetDaysSinceLastDrink();
    const today = new Date();
    const year = today.getFullYear();
    const month = today.getMonth() + 1; // Months are zero-based (0 = January)
    const day = today.getDate();
    let drinkBar = 0

    drinksToday = await GetDrinkCount(year,month,day)
    drinkBar = await GetDrinkTagColor(year,month,day);
    let tagText = await GetDrinks(year,month,day);
    tag = { text: tagText, color: tagColors[drinkBar] };
    progress = Math.min((drinkBar / 5) * 100, 100);
  });


</script>

<style>
  /* Txt colors */
  * {
    color: #76453B;
  }

  .container {
    padding: 20px;
    scrollbar-width: none;
    -ms-overflow-style: none;
  }
  
  .heading {
    text-align: center;
    font-size: 2rem;
    font-weight: bold;
    margin-bottom: 20px;
    color: #F8FAE5;
  }

  .form-container {
    display: flex;
    gap: 20px;
    justify-content: center;
    display: flex;
    width: 100%;
    justify-content: space-between;
    /* align-items: center; */
  }

  .form-box, .progress-box {
    background: #F8FAE5;
    border-radius: 10px;
    padding: 20px;
    box-shadow: 0 4px 6px rgba(118, 69, 59, 0.1);
  }

  .form-box {
    width: 60%;
    border: 1px solid #B19470;
    display: flex;
    justify-content: center;
  }

  /* .form-container {
      display: flex;
      width: 100%;
      justify-content: space-between;
      align-items: center;
  } */

  /* Left side takes up 70% */
  .form-left {
      flex: 6;
      display: flex;
      flex-direction: column;
      gap: 12px; /* Adds spacing between fields */
  }

  /* Right side takes up 30% */
  .form-right {
      flex: 4;
      display: flex;
      justify-content: center;
      align-items: center;
  }

  /* Floating divider between both sections */
  .divider {
      width: 2px;
      height: 50%; 
      background-color: #B19470;
      align-self: center;
  }

  /* Label & Input Field Styling */
  .form-group {
      display: flex;
      align-items: center;
  }

  .form-group label {
      color: #76453B;
      font-weight: bold;
      white-space: nowrap;
      margin-right: 10px;
  }

  input, select {
      padding: 6px;
      border: 1px solid #B19470;
      border-radius: 5px;
      font-size: 14px;
      text-align: center;
  }

  /* Fixed input sizes */
  .amount-cost-input {
      width: 55px;
  }

  .small-date {
      width: 35px;
  }

  .small-year {
      width: 50px;
  }

  /* Center the button in the right section */
  .styled-button {
      background: #B19470;
      color: white;
      border: none;
      padding: 10px 20px;
      border-radius: 5px;
      font-size: 16px;
      cursor: pointer;
      transition: background 0.3s ease;
  }

  .styled-button:hover {
      background: #8F6A4F;
  }

  /* .alcohol-input {
    appearance: auto;
    z-index: 1000;
    position: relative;
  } */







  /* TAB STUFF */

  .tab-buttons {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }

  .tab-button {
    padding: 10px 20px;
    margin: 0 5px;
    border: none;
    cursor: pointer;
    font-size: 1rem;
    border-radius: 5px;
    background: #B19470;
    color: #F8FAE5;
  }

  .tab-button.active {
    background: #F8FAE5;
    color: #76453B;
  }

  .tab-content {
    margin-top: 20px;
    padding: 20px;
    border: 1px solid #76453B;
    border-radius: 10px;
    background: #F8FAE5;
  }


  /*CALENDAR */
  .calendar-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
  }

  .year-controls {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    margin-bottom: 10px;
    width: 100%;
  }

  .arrow {
    cursor: pointer;
    padding: 5px 15px;
    font-size: 24px;
    user-select: none;
  }

  .calendar {
    display: grid;
    grid-template-columns: repeat(32, minmax(20px, 1fr)); /* Dynamic scaling */
    gap: 4px;
    width: 95vw; /* Ensures it stays responsive */
    max-width: 100%;
    overflow-x: auto;
    padding: 10px;
  }

  .cell {
    width: 100%;
    aspect-ratio: 1; /* Keeps cells square */
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: clamp(10px, 1vw, 14px); /* Responsive font size */
    border-radius: 5px;
  }

  .header {
    font-weight: bold;
  }

  .drinks-empty {
    background-color: gray;
    color: white;
  }

  .drinks-low {
    background-color: #60aa9b;
    color: white;
  }

  .drinks-moderate {
    background-color: #43766c;
    color: white;
  }

  .drinks-heavy {
    background-color: #ffdf60;
    color: white;
  }

  .drinks-binge {
    background-color: #fa8072;
    color: white;
  }

  .drinks-excessive {
    background-color: #ed4d09;
    color: white;
  }

  /* Progress Bar */
  .progress-box {
    display: flex;
    justify-content: space-between;
    /* align-items:start; */
    gap: 20px;
    width: 40%;
    border: 1px solid #B19470;
    display: flex;
  }

  .progress-left,
  .progress-right {
    flex: 1;
    padding: 10px;
  }

  .divider-progress {
    width: 2px;
    height: 80%; /* Adjust height to make it floating */
    background-color: #76453B; /* Line color */
    align-self: center; /* Keeps it vertically centered */
  }

  .counter {
    font-size: 1rem;
    color: #76453B;
    font-weight: bold;
    line-height: 1; /* Removes extra vertical spacing */
    margin-top: 0.7rem; /* Reduce bottom spacing */
    margin-bottom: 0.7rem;
  }

  .big-number {
    font-size: 3rem;
    font-weight: bold;
    color: #43766c;
    line-height: 1; /* Removes extra vertical spacing */
    margin-top: 0.7rem; /* Reduce bottom spacing */
    margin-bottom: 0.7rem;
  }

  .progress-container {
    width: 100%;
    margin-top: 8px;
  }

  .progress-bar {
    height: 10px;
    border-radius: 10px;
    background: #ddd;
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    border-radius: 10px;
    transition: width 0.3s ease-in-out;
  }

  .tag {
    display: inline-block;
    padding: 4px 8px;
    border-radius: 6px;
    font-size: 0.8rem;
    font-weight: bold;
    margin-top: 8px;
    color: white;
  }

  

</style>

<div class="container" style="--wails-draggable:drag">
  <!-- MODAL -->
  <Modal 
  isVisible={showModal} 
  onClose={closeModal}
  initialYear={selectedYear}
  initialMonth={selectedMonth}
  initialDay={selectedDay}
  onModalClose={onModalCloseCallback}
  />

  <h1 class="heading">Alcohol Tracker</h1>

  <div class="form-container">
    <!-- Form Box -->
    <div class="form-box">
      <div class="form-container">
        <!-- Left Side: Input Fields -->
        <div class="form-left">
          <div class="form-group">
            <label for="category">Alcohol (type)</label>
            <select id="category" bind:value={category} class="alcohol-input">
              <option value="" disabled selected>Alcohol</option>
              {#each alcoholCategories as type}
                <option value={type}>{type}</option>
              {/each}
            </select>
          </div>
    
          <div class="form-group">
            <label for="quantity">Amount (mL)</label>
            <input id="quantity" type="number" bind:value={quantity} class="amount-cost-input" />
          </div>
    
          <div class="form-group">
            <label for="cost">Cost ($)</label>
            <input id="cost" type="number" bind:value={cost} step="0.01" class="amount-cost-input" />
          </div>
    
          <div class="form-group date-group">
            <label>Date</label>
            <input id="day" type="number" bind:value={day} min="1" max="31" class="small-date" />
            <span>/</span>
            <input id="month" type="number" bind:value={month} min="1" max="12" class="small-date" />
            <span>/</span>
            <input id="year" type="number" bind:value={year} min="2000" max="2100" class="small-year" />
          </div>
        </div>
    
        <!-- Floating Divider -->
        <div class="divider"></div>
    
        <!-- Right Side: Add Entry Button -->
        <div class="form-right">
          <button class="styled-button" on:click={addEntry}>Add Entry</button>
        </div>
      </div>
    </div>
    
    
    
    

    <!-- Progress Box (Empty for Now) -->
    <div class="progress-box">
      <!-- Left Section -->
      <div class="progress-left">
        <p class="counter">Days Since Last Drink</p>
        <p class="big-number">{daysDisplay}</p>
      </div>
    
      <!-- Divider -->
      <div class="divider-progress"></div>
    
      <!-- Right Section -->
      <div class="progress-right">
        <p class="counter">Drinks Today</p>
        <p class="big-number">{drinksToday.toFixed(1)}</p>
    
        <div class="progress-container">
          <div class="progress-bar">
            <div class="progress-fill" style="width: {progress}%; background: {tag.color};"></div>
          </div>
        </div>
    
        <div class="tag" style="background: {tag.color};">{tag.text}</div>
      </div>
    </div>    
  </div>

  <!-- Tab Buttons -->
  <div class="tab-buttons">
    <button class="tab-button {activeTab === 'calendar' ? 'active' : ''}" on:click={() => setTab('calendar')}>Calendar</button>
    <button class="tab-button {activeTab === 'stats' ? 'active' : ''}" on:click={() => setTab('stats')}>Stats</button>
  </div>

  <!-- Tab Content -->
  <div class="tab-content" style="display: {activeTab === 'calendar' ? 'block' : 'none'};">
    <div class="calendar-container">
      <div class="year-controls">
        <span class="arrow" on:click={() => changeYear(-1)}>←</span>
        <span>{currentYear}</span>
        <span class="arrow" on:click={() => changeYear(1)}>→</span>
      </div>
    </div>
    
    <!-- Calendar Grid -->
    <div class="calendar">
      <!-- First row (Header: Days 1–31) -->
      <div class="cell"></div> <!-- Empty Top-Left -->
      {#each Array(31).fill(0).map((_, i) => i + 1) as day}
        <div class="cell header">{day}</div>
      {/each}
    
      <!-- Month Rows -->
      {#each months as month,monthIndex}
        <div class="cell header">{month}</div> <!-- Month Name -->
        {#each Array(31).fill(0).map((_, i) => i + 1) as day}
          {#if day <= daysInMonth[month]}
            {#await calendarDrinksMap(currentYear,monthIndex+1,day) then drinkClass}
              <div class="cell {drinkClass}"
              data-year="{currentYear}" 
              data-month="{monthIndex+1}" 
              data-day="{day}"
              tabindex="0"
              role="button"
              on:click={handleCalendarButton}
              on:keydown={handleCalendarButton}
              ></div>
            {/await}
          {:else}
            <div class="cell"></div>
          {/if}
        {/each}
      {/each}
    </div>    
  </div>

  <div class="tab-content" style="display: {activeTab === 'stats' ? 'block' : 'none'};">
    <h2>Stats View</h2>
    <p>Coming soon...</p>
  </div>
  <footer class="footer" style="margin-top: 1rem; color:#F8FAE5;">
    © {year} Akhil Sundaram. Built with <a href="https://wails.io" target="_blank">Wails</a>
  </footer>
</div>
