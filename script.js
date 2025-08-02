document.addEventListener('DOMContentLoaded', () => {
    const roomsList = document.getElementById('rooms-list');
    const bookingsList = document.getElementById('bookings-list');
    const bookingForm = document.getElementById('booking-form');
    const notification = document.getElementById('notification');
    
    // Modal elements
    const bookingModal = document.getElementById('booking-modal');
    const modalCloseBtn = document.getElementById('modal-close-btn');
    const modalRoomType = document.getElementById('modal-room-type');
    const modalRoomPrice = document.getElementById('modal-room-price');
    const modalRoomIdInput = document.getElementById('modal-room-id');

    const API_URL = 'http://localhost:8000';

    function showNotification(message, isError = false) {
        notification.textContent = message;
        notification.className = `fixed bottom-5 right-5 text-white py-2 px-4 rounded-lg shadow-xl transition-transform duration-500 ${isError ? 'bg-red-600' : 'bg-green-600'}`;
        notification.classList.remove('translate-x-[120%]');
        setTimeout(() => {
            notification.classList.add('translate-x-[120%]');
        }, 3000);
    }

    async function fetchRooms() {
        try {
            const response = await fetch(`${API_URL}/rooms`);
            if (!response.ok) throw new Error('Network response was not ok');
            const rooms = await response.json();
            roomsList.innerHTML = '';

            if (!rooms) return;

            rooms.forEach(room => {
                const isAvailable = room.available;
                const card = document.createElement('div');
                card.className = 'room-card bg-white rounded-lg shadow-md overflow-hidden transform transition-all duration-300 hover:shadow-xl hover:-translate-y-1';
                if (isAvailable) {
                    card.classList.add('cursor-pointer');
                } else {
                    card.classList.add('opacity-60', 'cursor-not-allowed');
                }
                
                card.innerHTML = `
                    <div class="relative h-56 overflow-hidden">
                        <img src="${room.imageUrl}" alt="${room.type}" class="room-image w-full h-full object-cover transition-transform duration-300">
                        <div class="absolute top-3 right-3 text-xs font-bold px-3 py-1 rounded-full ${isAvailable ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
                            ${isAvailable ? 'Available' : 'Booked'}
                        </div>
                    </div>
                    <div class="p-5">
                        <h3 class="text-xl font-serif">${room.type}</h3>
                        <p class="text-gray-600">Room ${room.id}</p>
                        <p class="text-2xl font-bold mt-3 text-blue-900">$${room.price.toFixed(2)}<span class="font-normal text-base text-gray-500">/night</span></p>
                    </div>
                `;

                if (isAvailable) {
                    card.addEventListener('click', () => openBookingModal(room));
                }
                roomsList.appendChild(card);
            });
        } catch (error) {
            console.error('Failed to fetch rooms:', error);
            roomsList.innerHTML = '<p class="text-red-500 col-span-full">Could not load rooms. Please ensure the server is running.</p>';
        }
    }

    function openBookingModal(room) {
        modalRoomType.textContent = room.type;
        modalRoomPrice.textContent = `Room ${room.id} - $${room.price.toFixed(2)} per night`;
        modalRoomIdInput.value = room.id;
        bookingModal.classList.remove('hidden'); // CORRECTED
    }

    function closeBookingModal() {
        bookingModal.classList.add('hidden'); // CORRECTED
        bookingForm.reset();
    }

    modalCloseBtn.addEventListener('click', closeBookingModal);
    bookingModal.addEventListener('click', (e) => {
        if (e.target === bookingModal) {
            closeBookingModal();
        }
    });


    async function fetchBookings() {
        try {
            const response = await fetch(`${API_URL}/bookings`);
            if (!response.ok) throw new Error('Network response was not ok');
            const bookings = await response.json();
            bookingsList.innerHTML = '';

            if (!bookings || bookings.length === 0) {
                bookingsList.innerHTML = '<p class="text-gray-500 text-center py-4">No active bookings.</p>';
                return;
            }

            bookings.forEach(booking => {
                const item = document.createElement('div');
                item.className = 'p-4 bg-blue-50 rounded-lg border border-blue-100 flex justify-between items-center';
                item.innerHTML = `
                    <div>
                        <p class="font-bold text-blue-900">Room ${booking.room}</p>
                        <p class="text-sm text-gray-600">${booking.guestName}</p>
                        <p class="text-xs text-gray-500 mt-1">${booking.checkIn} to ${booking.checkOut}</p>
                    </div>
                    <button class="delete-btn text-gray-400 hover:text-red-600 transition-colors" data-id="${booking._id}" title="Cancel Booking">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                    </button>
                `;
                bookingsList.appendChild(item);
            });
        } catch (error) {
            console.error('Failed to fetch bookings:', error);
            bookingsList.innerHTML = '<p class="text-red-500">Could not load bookings.</p>';
        }
    }

    bookingForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const newBooking = {
            room: modalRoomIdInput.value,
            guestName: bookingForm.guestName.value,
            checkIn: bookingForm.checkIn.value,
            checkOut: bookingForm.checkOut.value,
        };

        try {
            const response = await fetch(`${API_URL}/bookings`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(newBooking),
            });

            if (response.ok) {
                showNotification('Booking created successfully!');
                closeBookingModal();
                init();
            } else {
                showNotification('Failed to create booking. The room may have just been taken.', true);
            }
        } catch (error) {
            showNotification('An error occurred. Please try again.', true);
        }
    });

    bookingsList.addEventListener('click', async (e) => {
        const deleteButton = e.target.closest('.delete-btn');
        if (deleteButton) {
            const bookingId = deleteButton.dataset.id;
            if (!confirm('Are you sure you want to cancel this booking?')) return;

            try {
                const response = await fetch(`${API_URL}/bookings/${bookingId}`, { method: 'DELETE' });
                if (response.ok) {
                    showNotification('Booking cancelled.');
                    init();
                } else {
                    showNotification('Failed to cancel booking.', true);
                }
            } catch (error) {
                showNotification('An error occurred.', true);
            }
        }
    });

    function init() {
        fetchRooms();
        fetchBookings();
    }

    init();
});
