<script>
    export let isVisible = false;
    export let onClose = () => {};
    export let onModalClose = () => {};
  
    import { GetEntriesOnDate, DeleteDrink, GetAlcoholCategories, AddTrackerEntryUpdate, GetDrinkCount } from "../wailsjs/go/main/App";
    import { onMount } from "svelte";
    import { MdDeleteForever, MdEdit, MdCheck, MdClose } from "svelte-icons/md";
  
    export let initialYear = 2024;
    export let initialMonth = 1;
    export let initialDay = 1;
  
    let year = initialYear;
    let month = initialMonth;
    let day = initialDay;
    let drinkCount = 0.0;
    let entries = [];
    let editingIndex = null;
    let editAlcohol = "";
    let editQuantity = 1;
    let editCost = 1.0;
    let alcoholCategories = [];
  
    $: if (isVisible) {
        year = initialYear;
        month = initialMonth;
        day = initialDay;
        fetchEntries(year, month, day);
        getAlcoholTypes();
        getAlcoholDrinks(year,month,day);
    }

    async function getAlcoholTypes() {
        try {
        alcoholCategories = await GetAlcoholCategories();
        } catch (err) {
        console.error("Error fetching alcohol categories:", err);
        }
    }

    async function getAlcoholDrinks(year, month, day) {
        drinkCount = await GetDrinkCount(year, month,day);
    }
    
    async function fetchEntries(year, month, day) {
        try {
            const response = await GetEntriesOnDate(year, month, day);
            let parsedEntries = [];
            for (let key in response) {
                parsedEntries = parsedEntries.concat(response[key]);
            }
            entries = parsedEntries;
        } catch (error) {
            console.error("Error fetching entries:", error);
        }
    }
  
    function closeModal() {
        entries = [];
        editingIndex = null;
        onClose();
    }
  
    async function deleteEntry(year, month, day, alcohol, timestamp) {
        try {
            const deleteComplete = await DeleteDrink(year, month, day, alcohol, timestamp);
            if (deleteComplete) {
                entries = entries.filter(entry => !(entry.alcohol === alcohol && entry.timestamp === timestamp));
                onModalClose();
            }
        } catch (error) {
            console.error("Error deleting entry:", error);
        }
    }
  
    function modifyEntry(index) {
        editingIndex = index;
        editAlcohol = entries[index].alcohol;
        editQuantity = entries[index].quantity;
        editCost = entries[index].cost;


    }
  
    async function saveEntry(index) {
        try {
            const deleteComplete = await DeleteDrink(year, month, day, entries[index].alcohol, entries[index].timestamp);
            if (deleteComplete) {
                await AddTrackerEntryUpdate(year,month,day,editAlcohol,editQuantity,editCost, entries[index].timestamp)
                // entries = entries.filter(entry => !(entry.alcohol === alcohol && entry.timestamp === timestamp));
                onModalClose();
                entries[index].alcohol = editAlcohol;
                entries[index].quantity = editQuantity;
                entries[index].cost = editCost;
                editingIndex = null;
            }
        } catch (error) {
            console.error("Error deleting entry:", error);
        }
        
    }
  
    function cancelEdit() {
        editingIndex = null;
    }
  
    onMount(() => {
        if (isVisible) {
            fetchEntries(year, month, day);
        }
    });

    function handleKeyDown(event, action) {
      if (event.key === 'Enter' || event.key === ' ') {
          action();
      }
    }

    function validateInput(event, variable) {
        let value = event.target.value.trim();

        // Allow only positive numbers (integers or decimals, never zero or negative)
        if (!/^\d*\.?\d+$/.test(value) || Number(value) <= 0) {
        value = "1"; // Reset invalid input to 1 (ensuring > 0)
        }

        // Assign the validated value back to the correct variable
        if (variable === "cost") {
            editCost = parseFloat(value) || 1;
        } else if (variable === "ml") {
            editQuantity = parseFloat(value) || 1;
        }
    }
  </script>
  
  {#if isVisible}
    <div class="modal-overlay" on:click={closeModal} on:keydown={event => handleKeyDown(event, closeModal)}>
        <div class="modal-content" on:click|stopPropagation on:keydown={event => handleKeyDown(event, closeModal)}>
            <h3>{day}/{month}/{year}</h3>
            <!-- <h2>drinks: {drinkCount.toFixed(1)}</h2> -->
            
            {#if entries.length > 0}
                <div class="entries-container">
                    {#each entries as entry, index}
                        <div class="entry-card">
                            {#if editingIndex === index}
                            <div class="entry-info-container editing">
                                <div class="edit-input-container">
                                    <select class="alcohol-type" bind:value={editAlcohol}>
                                        {#each alcoholCategories as category}
                                            <option value={category}>{category}</option>
                                        {/each}
                                    </select>                            
                                    <input class="entry-info" type="number" bind:value={editQuantity} min="0" step="1" on:input={(e) => validateInput(e,"ml")} />
                                    <input class="entry-info" type="number" bind:value={editCost} min="0" step="0.01" on:input={(e) => validateInput(e,"cost")} />
                                </div>
                            
                                <div class="divider"></div>
                            
                                <div class="edit-action-buttons">
                                    <button class="save-button" on:click={() => saveEntry(index)}>
                                        <MdCheck size="24" color="#28a745" />
                                    </button>
                                    <button class="cancel-button" on:click={cancelEdit}>
                                        <MdClose size="24" color="#dc3545" />
                                    </button>
                                </div>
                            </div>
                            
                            {:else}
                                <div class="entry-info-container">
                                    <div>
                                        <p class="alcohol-type">{entry.alcohol}</p>
                                        <p class="entry-info">Cost: ${entry.cost.toFixed(2)}</p>
                                        <p class="entry-info">Quantity: {entry.quantity} mL</p>
                                    </div>
                                    <div class="action-buttons">
                                        <button class="modify-button" on:click={() => modifyEntry(index)}>
                                            <MdEdit size="24" color="#007bff" />
                                        </button>
                                        <button class="delete-button" on:click={() => deleteEntry(year, month, day, entry.alcohol, entry.timestamp)}>
                                            <MdDeleteForever size="24" color="#dc3545" />
                                        </button>
                                    </div>
                                </div>
                            {/if}
                        </div>
                    {/each}
                </div>
            {:else}
                <p>We call this day, A Liver Success!</p>
            {/if}
            <button on:click={closeModal}>Close</button>
        </div>
    </div>
  {/if}

<style>
  .modal-overlay {
      position: fixed;
      top: 0;
      left: 0;
      width: 100vw;
      height: 100vh;
      background: rgba(0, 0, 0, 0.5);
      border: 1.5px solid #4C4B16;
      backdrop-filter: blur(5px);
      display: flex;
      justify-content: center;
      align-items: center;

  }

  .modal-content {
      background: #B19470;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
      border: 1.5px solid #4C4B16;
      z-index: 100;
      text-align: center;
      width: 300px;
      max-height: 80vh;
      overflow-y: auto;
  }

  .entries-container {
      display: flex;
      flex-direction: column;
      gap: 10px;
      margin-top: 15px;
  }

  .entry-card {
      background: #F8FAE5;
      padding: 10px;
      border-radius: 8px;
      box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.1);
      position: relative;
      display: flex;
      align-items: center;
      justify-content: space-between;
  }

  .entry-info-container {
      display: flex;
      justify-content: space-between;
      align-items: center;
      width: 100%;
      position: relative;
      flex-wrap: wrap;
  }

  .alcohol-type {
      font-weight: bold;
      font-size: 16px;
      margin-bottom: 5px;
  }

  .entry-info {
      font-size: 14px;
      /* color: #; */
      margin: 2px 0;
  }

  /* Action Buttons */
  .action-buttons {
      position: absolute;
      top: 10px;
      right: 10px;
      display: flex;
      gap: 5px; 
  }

  .modify-button,
  .delete-button,
  .save-button,
  .cancel-button {
      background: #76453B;
      color: #F8FAE5;
      border: none;
      border-radius: 50%;
      width: 24px;
      height: 24px;
      font-size: 14px;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 2px;
  }

  .modify-button:hover {
      background: #FFEAC5;
      color: #76453B;
  }

  .delete-button:hover {
      background: #ed4d09;
      color: #F8FAE5;
  }

  button {
      margin-top: 15px;
      padding: 8px 12px;
      border: none;
      background: #76453B;
      color: #F8FAE5;
      border-radius: 5px;
      cursor: pointer;
  }

  button:hover {
      background: #B17457;
      color: #F8FAE5;
  }


  .entry-info-container input {
      width: 70px;
      padding: 5px;
      margin: 2px;
      border: 1px solid #76453B;
      border-radius: 4px;
      color: #76453B;
  }
  .save-button, .cancel-button {
    color: #F8FAE5;
    background: #76453B;
    border: none;
    cursor: pointer;
  }

  .cancel-button:hover {
      background-color: #f39360;
      /* color: #FFEAC5; */
  }

  .save-button:hover {
      background-color: #B5CFB7;
  }

  .entry-info-container.editing {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    position: relative;
  }

    .entry-info-container input {
        flex-grow: 1; /* Ensures inputs take available space */
        padding: 5px;
        min-width: 70px; /* Prevents inputs from shrinking too much */
        margin-right: 10px; /* Adds spacing between inputs */
    }

    .entry-info-container .edit-action-buttons {
        display: flex;
        gap: 5px;
        align-items: center;
        justify-content: flex-end;
        min-width: 80px; /* Ensures the buttons don't get too close */
    }

    .edit-input-container {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;
        gap: 10px;
    }

    .divider {
        width: 75%;
        height: 1px;
        background: #ccc;
        margin: 10px auto;
    }

    .entry-info-container {
        display: flex;
        flex-direction: column; /* Stack elements */
        align-items: center; /* Center elements */
        justify-content: center; /* Ensure vertical centering */
        width: 100%;
    }

    .save-button,
    .cancel-button {
        width: 26px; /* Increase button width */
        height: 26px; /* Increase button height */
        padding: 3px;
        margin-left: 1rem;
        margin-right: 1rem;
        margin-top: 0.2rem;
    }

    select {
        color: #76453B;
    }






</style>
