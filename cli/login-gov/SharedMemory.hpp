//
// Created by Nicholas Robison on 2019-01-29.
//

#ifndef CMSAUTHZCLI_SHAREDMEMORY_HPP
#define CMSAUTHZCLI_SHAREDMEMORY_HPP

#include <boost/interprocess/sync/interprocess_semaphore.hpp>

template<typename Alloc>
struct SharedMemory {

    using allocator_type = typename Alloc::template rebind<char>::other;

public:

    explicit SharedMemory(size_t aSize, Alloc alloc = {}) : mutex(0), done(0), mData(aSize, alloc) {
        // Not used
    }

    std::vector<char, Alloc> mData;
    boost::interprocess::interprocess_semaphore mutex, done;
};

#endif //CMSAUTHZCLI_SHAREDMEMORY_HPP
